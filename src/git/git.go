package git

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const (
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
)

var pattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
	Warn(s string, args ...interface{})
	Debug(s string, args ...interface{})
}

type GitMethods struct {
	getBranchPointedToHead func() (*plumbing.Reference, error)
	getCommitHistory       func() ([]*object.Commit, error)
	getMostRecentCommit    func() (CommitInfo, error)
	getAllTags             func() ([]object.Tag, error)
	getMostRecentTag       func() (string, error)
	addToStage             func() error
	commitChanges          func(newReleaseVersion string) error
	push                   func() error
	tagExists              func(tag string) (bool, error)
	setTag                 func(tag string) error
	pushTags               func() error
}

type ElapsedTime func(functionName string) func()

type GitVersioning struct {
	git                  GitMethods
	log                  Logger
	printElapsedTime     ElapsedTime
	url                  string
	destinationDirectory string
	username             string
	password             string
	repo                 *git.Repository
	branchHead           *plumbing.Reference
	commitHistory        []*object.Commit
	tagsList             []object.Tag
	mostRecentCommit     CommitInfo
	mostRecentTag        string
}

type CommitInfo struct {
	Hash        string
	AuthorName  string
	AuthorEmail string
	Message     string
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (g *GitVersioning) validate() error {
	if g.url == "" {
		return errors.New("url cannot be empty")
	}

	if g.destinationDirectory == "" {
		return errors.New("destination directory cannot be empty")
	}

	if g.username == "" {
		return errors.New("username cannot be empty")
	}

	if g.password == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

func (g *GitVersioning) GetChangeHash() string {
	return g.mostRecentCommit.Hash
}

func (g *GitVersioning) GetChangeAuthorName() string {
	return g.mostRecentCommit.AuthorName
}

func (g *GitVersioning) GetChangeAuthorEmail() string {
	return g.mostRecentCommit.AuthorEmail
}

func (g *GitVersioning) GetChangeMessage() string {
	return g.mostRecentCommit.Message
}

func (g *GitVersioning) GetCurrentVersion() string {
	return g.mostRecentTag
}

func (g *GitVersioning) UpgradeRemoteRepository(newVersion string) error {
	// consider 1.0.0 as the start tag of a repository when it does not have tags yet
	if newVersion == "0.1.0" || newVersion == "0.0.1" {
		newVersion = "1.0.0"
	}

	if err := g.git.commitChanges(newVersion); err != nil {
		return fmt.Errorf("error during commit operation due to: %w", err)
	}

	if err := g.git.push(); err != nil {
		return fmt.Errorf("error during push operation due to: %w", err)
	}

	if err := g.git.setTag(newVersion); err != nil {
		return fmt.Errorf("error during set tag operation due to: %w", err)
	}

	if err := g.git.pushTags(); err != nil {
		return fmt.Errorf("error during push tags operation due to: %w", err)
	}

	return nil
}

func (g *GitVersioning) getBranchPointedToHead() (*plumbing.Reference, error) {
	defer g.printElapsedTime("GetBranchPointedToHead")()
	g.log.Info("getting branch pointed to HEAD")
	ref, err := g.repo.Head()
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func (g *GitVersioning) getCommitHistory() ([]*object.Commit, error) {
	defer g.printElapsedTime("GetComitHistory")()
	g.log.Info("getting commit history")
	cIter, err := g.repo.Log(&git.LogOptions{From: g.branchHead.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func (g *GitVersioning) isTimeAfter(timeToCheck, referenceTime time.Time) bool {
	return timeToCheck.After(referenceTime)
}

func (g *GitVersioning) getMostRecentCommit() (CommitInfo, error) {

	if len(g.commitHistory) == 0 {
		return CommitInfo{}, fmt.Errorf("no commits found")
	}

	recentCommit := g.commitHistory[0]
	for _, commit := range g.commitHistory {
		if g.isTimeAfter(commit.Author.When, recentCommit.Author.When) {
			recentCommit = commit
		}
	}

	g.log.Debug("Most recent commit Author Name: ", recentCommit.Author.Name)
	g.log.Debug("Most recent commit Author Email: ", recentCommit.Author.Email)
	g.log.Debug("Most recent commit Message: ", recentCommit.Message)
	g.log.Debug("Most recent commit Time: ", recentCommit.Author.When)

	result := CommitInfo{
		Hash:        recentCommit.Hash.String(),
		AuthorName:  recentCommit.Author.Name,
		AuthorEmail: recentCommit.Author.Email,
		Message:     recentCommit.Message,
	}

	return result, nil
}

func (g *GitVersioning) getAllTags() ([]object.Tag, error) {
	defer g.printElapsedTime("getAllTags")()
	g.log.Info("getting all tags from repository")

	tagsIter, err := g.repo.Tags()

	if err != nil {
		return nil, err
	}

	var tags []object.Tag
	if err := tagsIter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, object.Tag{
			Hash: ref.Hash(),
			Name: ref.Name().String(),
		})

		return nil
	}); err != nil {
		return nil, err
	}

	return tags, nil
}

func (g *GitVersioning) getMostRecentTag() (string, error) {
	defer g.printElapsedTime("GetMostRecentTag")()
	g.log.Info("getting most recent tag from repository")

	if len(g.tagsList) == 0 {
		return "0.0.0", nil
	}

	mapTags := make(map[*Version]string)

	for _, currentTag := range g.tagsList {
		tag := strings.TrimSpace(strings.Replace(currentTag.Name, "refs/tags/", "", 1))

		if pattern.MatchString(tag) {
			version := newVersion(tag)
			mapTags[version] = tag
		}
	}

	var latest *Version
	var latestTag string
	for version, tag := range mapTags {
		latest, latestTag = isSetNewVersion(latest, version, tag)
	}

	if latestTag == "" {
		return "0.0.0", nil
	}

	return latestTag, nil
}

func (g *GitVersioning) addToStage() error {
	worktree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}
	g.log.Info(colorGreen + "Changes added to stage area..." + colorReset)
	return nil
}

func (g *GitVersioning) commitChanges(newReleaseVersion string) error {
	if err := g.git.addToStage(); err != nil {
		return err
	}

	worktree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	signature := &object.Signature{Name: g.mostRecentCommit.AuthorName, Email: g.mostRecentCommit.AuthorEmail, When: time.Now()}

	message := fmt.Sprintf("type: [skip]: message: Commit automatically generated by Semantic Release. The new tag is %s", newReleaseVersion)
	commit, err := worktree.Commit(message, &git.CommitOptions{Author: signature, Committer: signature})
	if err != nil {
		return err
	}

	g.log.Info(colorGreen+"New commit added: %s"+colorReset, commit.String())
	return nil
}

func (g *GitVersioning) push() error {
	err := g.repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: g.username,
			Password: g.password,
		},
		InsecureSkipTLS: true})
	if err != nil {
		return err
	}

	return nil
}

func (g *GitVersioning) tagExists(tag string) (bool, error) {
	res := false
	tags, err := g.repo.TagObjects()
	if err != nil {
		return res, fmt.Errorf("error while getting tags due to %w", err)
	}

	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf("tag %s already exists", tag)
		}
		return nil
	})
	if err != nil && err.Error() != fmt.Sprintf("tag %s already exists", tag) {
		return false, fmt.Errorf("error while iterating tags due to: %w", err)
	}
	return res, nil
}

func (g *GitVersioning) setTag(tag string) error {
	g.log.Info("Set tag %s", tag)
	tagExists, err := g.git.tagExists(tag)
	if err != nil {
		return err
	}
	if tagExists {
		return err
	}

	g.log.Info("Creating tag %s", tag)
	err = g.setBranchHead()
	if err != nil {
		return err
	}
	_, err = g.repo.CreateTag(tag, g.branchHead.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  g.mostRecentCommit.AuthorName,
			Email: g.mostRecentCommit.AuthorEmail,
			When:  time.Now(),
		},
		Message: fmt.Sprintf("Generated by semantic-release %s", tag),
	})
	if err != nil {
		return fmt.Errorf("error while creating tag due to: %w", err)
	}
	g.log.Info("Tag %s successfully created", tag)
	return nil
}

func (g *GitVersioning) pushTags() error {

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth: &http.BasicAuth{
			Username: g.username,
			Password: g.password,
		},
		InsecureSkipTLS: true,
	}
	err := g.repo.Push(po)

	if err != nil {
		return err
	}

	return nil
}

func (g *GitVersioning) cloneRepoToDirectory() (*git.Repository, error) {
	defer g.printElapsedTime("CloneRepoToDirectory")()

	g.log.Info(colorYellow+"cloning repo "+colorCyan+" %s "+colorYellow+" into "+colorCyan+"%s"+colorReset, g.url, g.destinationDirectory)
	repo, err := git.PlainClone(g.destinationDirectory, false, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      g.url,
		Auth: &http.BasicAuth{Username: g.username,
			Password: g.password,
		},
		InsecureSkipTLS: true,
	})

	if err == nil {
		return repo, nil
	}

	if err == git.ErrRepositoryAlreadyExists {
		g.log.Warn("repository was already cloned")
		return git.PlainOpen(g.destinationDirectory)
	}
	g.log.Error("error while cloning gitab repository due to: %s", err)
	return nil, err
}

func (g *GitVersioning) setBranchHead() error {
	branchHead, err := g.git.getBranchPointedToHead()
	if err != nil {
		return fmt.Errorf("error while retrieving the branch pointed to HEAD due to: %w", err)
	}
	g.branchHead = branchHead
	return nil
}

func (g *GitVersioning) initialize() error {
	err := g.setBranchHead()
	if err != nil {
		return err
	}

	commitHistory, err := g.git.getCommitHistory()
	if err != nil {
		return fmt.Errorf("error while retrieving the commit history  due to: %w", err)
	}
	g.commitHistory = commitHistory

	mostRecentCommit, err := g.git.getMostRecentCommit()
	if err != nil {
		return fmt.Errorf("error while retrieving tags from repository due to: %w", err)
	}
	g.mostRecentCommit = mostRecentCommit

	allTags, err := g.git.getAllTags()
	if err != nil {
		return fmt.Errorf("errow while getting all tags due to: %w", err)
	}
	g.tagsList = allTags

	mostRecentTag, err := g.git.getMostRecentTag()
	if err != nil {
		return fmt.Errorf("error while getting most recent tage due to: %w", err)
	}
	g.mostRecentTag = mostRecentTag

	return nil
}

func newVersion(tag string) *Version {
	segments := strings.Split(tag, ".")
	major, _ := strconv.Atoi(segments[0])
	minor, _ := strconv.Atoi(segments[1])
	patch, _ := strconv.Atoi(segments[2])

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func isSetNewVersion(latest, version *Version, tag string) (*Version, string) {
	if latest == nil || version.isGreaterThan(latest) {
		return version, tag
	}
	return latest, ""
}

func (v *Version) isGreaterThan(other *Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch > other.Patch
}

func New(log Logger, printElapsedTime ElapsedTime, url, username, password, destinationDirectory string) (*GitVersioning, error) {
	gitLabVersioning := &GitVersioning{
		log:                  log,
		printElapsedTime:     printElapsedTime,
		username:             username,
		password:             password,
		url:                  url,
		destinationDirectory: destinationDirectory,
	}

	if err := gitLabVersioning.validate(); err != nil {
		return nil, err
	}

	repo, err := gitLabVersioning.cloneRepoToDirectory()
	if err != nil {
		return nil, fmt.Errorf("error while initiating git package due to : %w", err)
	}

	gitLabVersioning.repo = repo

	gitLabVersioning.git = GitMethods{
		getBranchPointedToHead: gitLabVersioning.getBranchPointedToHead,
		getCommitHistory:       gitLabVersioning.getCommitHistory,
		getMostRecentCommit:    gitLabVersioning.getMostRecentCommit,
		getAllTags:             gitLabVersioning.getAllTags,
		getMostRecentTag:       gitLabVersioning.getMostRecentTag,
		addToStage:             gitLabVersioning.addToStage,
		commitChanges:          gitLabVersioning.commitChanges,
		push:                   gitLabVersioning.push,
		tagExists:              gitLabVersioning.tagExists,
		setTag:                 gitLabVersioning.setTag,
		pushTags:               gitLabVersioning.pushTags,
	}

	if err := gitLabVersioning.initialize(); err != nil {
		return nil, err
	}

	return gitLabVersioning, nil
}
