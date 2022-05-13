package git

import (
	"errors"
	"fmt"
	"os"
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

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
	Warn(s string, args ...interface{})
}

type ElapsedTime func(functionName string) func()

type GitVersioning struct {
	log                  Logger
	printElapsedTime     ElapsedTime
	url                  string
	destinationDirectory string
	repo                 *git.Repository
	username             string
	password             string
	branchHead           *plumbing.Reference
	commitHistory        []*object.Commit
	mostRecentCommit     commitInfo
	tagsList             []object.Tag
	mostRecentTag        string
}

type commitInfo struct {
	hash        string
	authorName  string
	authorEmail string
	message     string
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
	return g.mostRecentCommit.hash
}

func (g *GitVersioning) GetChangeAuthorName() string {
	return g.mostRecentCommit.authorName
}

func (g *GitVersioning) GetChangeAuthorEmail() string {
	return g.mostRecentCommit.authorEmail
}

func (g *GitVersioning) GetChangeMessage() string {
	return g.mostRecentCommit.message
}

func (g *GitVersioning) GetCurrentVersion() string {
	return g.mostRecentTag
}

func (g *GitVersioning) UpgradeRemoteRepository(newVersion string) error {
	// consider 1.0.0 as the start tag of a repository when it does not have tags yet
	if newVersion == "0.1.0" || newVersion == "0.0.1" {
		newVersion = "1.0.0"
	}

	if err := g.commitChanges(newVersion); err != nil {
		return err
	}

	if err := g.push(); err != nil {
		return err
	}

	if err := g.setTag(newVersion); err != nil {
		return err
	}

	if err := g.pushTags(); err != nil {
		return err
	}

	return nil
}

func (g *GitVersioning) getBranchPointedToHead() error {
	defer g.printElapsedTime("GetBranchPointedToHead")()
	g.log.Info("getting branch pointed to HEAD")
	ref, err := g.repo.Head()
	if err != nil {
		return fmt.Errorf("error while retrieving the branch pointed to HEAD due to: %w", err)
	}

	g.branchHead = ref

	return nil
}

func (g *GitVersioning) getCommitHistory() error {
	defer g.printElapsedTime("GetComitHistory")()
	g.log.Info("getting commit history")
	cIter, err := g.repo.Log(&git.LogOptions{From: g.branchHead.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return fmt.Errorf("error while retrieving the commit history  due to: %w", err)
	}

	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error while iterating over the commits due to: %w", err)
	}

	g.commitHistory = commits

	return nil
}

func (g *GitVersioning) getMostRecentCommit() error {

	if len(g.commitHistory) == 0 {
		return fmt.Errorf("no commits found")
	}

	recentCommit := g.commitHistory[0]
	for i, commit := range g.commitHistory {
		if i > 0 {
			if commit.Author.When.After(g.commitHistory[i-1].Author.When) {
				recentCommit = commit
			}
		}
	}
	fmt.Println("HI")
	fmt.Println(recentCommit)
	fmt.Println(recentCommit.Hash.String())
	g.mostRecentCommit = commitInfo{
		hash:        recentCommit.Hash.String(),
		authorName:  recentCommit.Author.Name,
		authorEmail: recentCommit.Author.Email,
		message:     recentCommit.Message,
	}
	return nil
}

func (g *GitVersioning) getAllTags() error {
	defer g.printElapsedTime("getAllTags")()
	g.log.Info("getting all tags from repository")

	tagsIter, err := g.repo.Tags()
	errMessage := "error while retrieving tags from repository due to: %w"
	if err != nil {
		return fmt.Errorf(errMessage, err)
	}

	var tags []object.Tag
	if err := tagsIter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, object.Tag{
			Hash: ref.Hash(),
			Name: ref.Name().String(),
		})

		if len(tags) == 0 {
			return fmt.Errorf(errMessage, "no tags found.")
		}
		return nil
	}); err != nil {
		return fmt.Errorf(errMessage, err)
	}

	g.tagsList = tags

	return nil
}

func (g *GitVersioning) getMostRecentTag() error {
	defer g.printElapsedTime("GetMostRecentTag")()
	g.log.Info("getting most recent tag from repository")

	if len(g.tagsList) == 0 {
		g.mostRecentTag = "0.0.0"
		return nil
	}

	mapTags := make(map[int]string)

	for _, currentTag := range g.tagsList {
		tag := strings.TrimSpace(strings.Replace(currentTag.Name, "refs/tags/", "", 1))

		tagOnlyNumbers := strings.ReplaceAll(tag, ".", "")
		tagInt, err := strconv.Atoi(tagOnlyNumbers)
		if err != nil {
			return fmt.Errorf("error while getting most recent tage due to: could not convert %v to int", tagOnlyNumbers)
		}
		mapTags[tagInt] = tag
	}

	previous := 0
	for key, element := range mapTags {
		if key > previous {
			previous = key
			g.mostRecentTag = element
		}
	}

	return nil
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
	if err := g.addToStage(); err != nil {
		return err
	}

	worktree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	signature := &object.Signature{Name: g.mostRecentCommit.authorName, Email: g.mostRecentCommit.authorEmail, When: time.Now()}

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
		return fmt.Errorf("error while pushing commits to remote repository due to: %w", err)
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
	tagExists, err := g.tagExists(tag)
	if err != nil {
		return err
	}
	if tagExists {
		return err
	}

	g.log.Info("Creating tag %s", tag)
	_, err = g.repo.CreateTag(tag, g.branchHead.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  g.mostRecentCommit.authorName,
			Email: g.mostRecentCommit.authorEmail,
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
		return fmt.Errorf("error while pushing tag to remote branch due to: %w", err)
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

func (g *GitVersioning) initialize() error {
	if err := g.getBranchPointedToHead(); err != nil {
		return err
	}

	if err := g.getCommitHistory(); err != nil {
		return err
	}

	if err := g.getMostRecentCommit(); err != nil {
		return err
	}

	if err := g.getAllTags(); err != nil {
		return err
	}

	if err := g.getMostRecentTag(); err != nil {
		return err
	}

	return nil
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
		gitLabVersioning.log.Error(err.Error())
		return nil, err
	}

	repo, err := gitLabVersioning.cloneRepoToDirectory()
	if err != nil {
		return nil, fmt.Errorf("error while initiating git package due to : %w", err)
	}

	gitLabVersioning.repo = repo

	if err := gitLabVersioning.initialize(); err != nil {
		return nil, err
	}

	return gitLabVersioning, nil
}
