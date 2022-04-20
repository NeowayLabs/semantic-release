package git

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	e "github.com/NeowayLabs/semantic-release/src/errors"
	style "github.com/NeowayLabs/semantic-release/src/style"
	timeutils "github.com/NeowayLabs/semantic-release/src/time"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type GitService interface {
	CloneRepoToDirectory(url, destinationDir string, keepCloned bool, auth *ssh.PublicKeys) (*git.Repository, error)
	GetBranchPointedToHead(repo *git.Repository) (*plumbing.Reference, error)
	GetCommitHistory(repo *git.Repository, ref *plumbing.Reference) ([]*object.Commit, error)
	GetMostRecentCommit(commits []interface{}) (*object.Commit, error)
	GetMostRecentTag(repo *git.Repository) (string, error)
	AddToStage(repo *git.Repository) error
	CommitChanges(repo *git.Repository, authorName string, authorEmail string, commitMessage string) error
	Push(repo *git.Repository, auth *ssh.PublicKeys) error
	SetTag(repo *git.Repository, tag, authorName, authorEmail string) (bool, error)
	DeleteTag(repo *git.Repository, tag string) (bool, error)
	PushTags(repo *git.Repository, auth *ssh.PublicKeys) error
}

type gitService struct{}

type Signature struct {
	Name  string
	Email string
	When  time.Time
}

// CloneRepoToDirectory aims to clone the repository from remote to local.
// Args:
// 		url (string): URL from where the repository must be cloned.
// 		destinationDir (string): Destination directory to of the repository.
//		keepCloned (bool): Defines if the repository must be cloned once again even if the destinationDir already exists.
// 		auth (*ssh.PublicKeys): Authorization key.
// Returns:
// 		*git.Repository: Returns a repository reference.
// 		err: Error whenever unexpected issues happen.
func (s *gitService) CloneRepoToDirectory(url, destinationDir string, keepCloned bool, auth *ssh.PublicKeys) (*git.Repository, error) {
	defer timeutils.GetElapsedTime("CloneRepoToDirectory")()
	log.Printf(style.Yellow+"cloning repo "+style.Cyan+" %s "+style.Yellow+" into "+style.Cyan+"%s"+style.Reset, url, destinationDir)
	repo, err := git.PlainClone(destinationDir, false, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      url,
		Auth:     auth,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			log.Println(e.ErrMsgRepoAlreadyCloned)
			if keepCloned {
				log.Printf("removing path %s", destinationDir)
				err := os.RemoveAll(destinationDir)
				if err != nil {
					return nil, err
				}
				return s.CloneRepoToDirectory(url, destinationDir, false, auth)
			}
			return git.PlainOpen(destinationDir)
		} else {
			log.Printf("clone git repo error: %s", err)
			return nil, err
		}
	}
	return repo, nil
}

// GetBranchPointedToHead aims to get the branch reference pointing to HEAD.
// Args:
// 		repo (*git.Repository): Repository to get the branch from.
// Returns:
// 		[]*object.Commit: Returns a branch reference.
// 		err: Error whenever unexpected issues happen.
func (s *gitService) GetBranchPointedToHead(repo *git.Repository) (*plumbing.Reference, error) {
	defer timeutils.GetElapsedTime("GetBranchPointedToHead")()
	log.Println("getting branch pointed to HEAD")
	ref, err := repo.Head()
	if e.Error(err, e.ErrMsgRetrievingBranchHead) {
		return nil, err
	}

	return ref, nil
}

// GetCommitHistory aims to retrieve the commit history from a given repository.
// Args:
// 		repo (*git.Repository): Repository to get the commits from.
// 		ref (*plumbing.Reference): Branch reference where the commits must be retrieved from.
// Returns:
// 		[]*object.Commit: Returns a slice of commits.
// 		err: Error whenever unexpected issues happen.
func (s *gitService) GetCommitHistory(repo *git.Repository, ref *plumbing.Reference) ([]*object.Commit, error) {
	defer timeutils.GetElapsedTime("GetComitHistory")()
	log.Println("getting commit history")
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
	if e.Error(err, e.ErrMsgRetrievingCommitHistory) {
		return nil, err
	}

	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	if e.Error(err, e.ErrMsgIteratingCommitHistory) {
		return nil, err
	}

	return commits, nil
}

// GetMostRecentCommit is responsible to get the first commit from a slice of commits sorted in ascending order
// Args:
// 		commits ([]*object.Commit): An input slice of commits sorted in ascending order
// Returns:
// 		*object.Commit: The most recent commit from a slice of commits sorted in ascending order.
func (s *gitService) GetMostRecentCommit(commits []interface{}) (*object.Commit, error) {

	if len(commits) == 0 {
		log.Println(e.ErrMsgNoCommitsFound)
		return nil, errors.New(e.ErrMsgNoCommitsFound)
	}

	commitList := []*object.Commit{}
	for _, commit := range commits {
		commitParsed, ok := commit.(*object.Commit)
		if !ok {
			return nil, errors.New(e.ErrMsgInterfaceToStruct)
		}
		commitList = append(commitList, commitParsed)
	}

	return commitList[0], nil
}

// getAllTags aims to retrieve all the tag from the repository.
// Args:
// 		repo (*git.Repository): Git repository where the tags must be found.
// Returns:
// 		[]object.Tag: Slice of object.Tags containing the tags information.
// 		err: Error whenever unexpected issues happen.
func (s *gitService) getAllTags(repo *git.Repository) ([]object.Tag, error) {
	defer timeutils.GetElapsedTime("getAllTags")()
	log.Println("getting all tags from repository")

	tagsIter, err := repo.Tags()
	if e.Error(err, e.ErrMsgRetrievingTags) {
		return nil, err
	}

	var tags []object.Tag
	if err := tagsIter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, object.Tag{
			Hash: ref.Hash(),
			Name: ref.Name().String(),
		})

		if len(tags) == 0 {
			return errors.New(e.ErrMsgRetrievingTags)
		}
		return nil
	}); err != nil {
		return nil, errors.New(e.ErrMsgRetrievingTags)
	}

	return tags, nil
}

// GetMostRecentTag aims to get the most recent tag from the repository.
// Args:
// 		repo (*git.Repository): Git repository where the tags must be found.
// Returns:
// 		string: Returns the most recent tag. Default to `0.0.0`
// 		err: Error whenever unexpected issues happen.
func (s *gitService) GetMostRecentTag(repo *git.Repository) (string, error) {
	defer timeutils.GetElapsedTime("GetMostRecentTag")()
	log.Println("getting most recent tag from repository")

	tags, err := s.getAllTags(repo)
	if err != nil {
		return "", errors.New(e.ErrMsgGettingMostRecentTag)
	}

	if len(tags) == 0 {
		return "0.0.0", nil
	}

	lastPosition := tags[len(tags)-1]

	result := strings.TrimSpace(strings.Replace(lastPosition.Name, "refs/tags/", "", 1))

	return result, nil
}

// setDefaultSignature aims to fillup and return a struct with the same fieds as object.Signature.
// Args:
// 		name (string): Git user name.
// 		email (string): Git user email.
// Returns:
// 		Signature: Returns a new Signature.
func setDefaultSignature(name, email string) Signature {
	return Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}
}

// addToStage aims to add all the local changes to the stage area.
// Args:
// 		repo (*git.Repository): Git repository that was changed.
// Returns:
// 		error: Returns an error whenever unexpected issues happen.
func (s *gitService) AddToStage(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}
	log.Println(style.Green + "Changes added to stage area..." + style.Reset)
	return nil
}

// CommitChanges aims to commit all the staging changes to the transfer area.
// Args:
// 		repo (*git.Repository): Git repository containing the changes.
// 		authorName (string): Author name.
// 		authorEmail (string): Author email.
// 		commitMessage (string): Commit message.
// Returns:
// 		error: Returns an error whenever unexpected issues happen.
func (s *gitService) CommitChanges(repo *git.Repository, authorName string, authorEmail string, commitMessage string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	author := setDefaultSignature(authorName, authorEmail)
	signature := &object.Signature{Name: author.Name, Email: authorEmail, When: author.When}
	commit, err := worktree.Commit(commitMessage, &git.CommitOptions{Author: signature, Committer: signature})
	if err != nil {
		return err
	}

	log.Printf(style.Green+"New commit added: %s"+style.Reset, commit.String())
	return nil
}

// CommitChanges aims to push the commits to the remote repository.
// Args:
// 		repo (*git.Repository): Git repository containing the changes.
// 		auth (*ssh.PublicKeys): SSH key to authenticate while pushing the commits to the remote repository.
// Returns:
// 		error: Returns an error whenever unexpected issues happen.
func (s *gitService) Push(repo *git.Repository, auth *ssh.PublicKeys) error {
	err := repo.Push(&git.PushOptions{Auth: auth})
	if err != nil {
		return err
	}

	return nil
}

// tagExists aims to find a given tag in the repository.
// Args:
// 		repo (*git.Repository): Git repository where the tag must be found.
// 		tag (string): Tag name to search.
// Returns:
// 		bool: True when tag found, otherwise false.
func tagExists(repo *git.Repository, tag string) bool {
	tags, err := repo.TagObjects()
	if err != nil {
		log.Fatalf(e.ErrMsgGetingTags, err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(e.ErrMsgTagExists, tag)
		}
		return nil
	})
	if err != nil && err.Error() != fmt.Sprintf(e.ErrMsgTagExists, tag) {
		log.Fatalf(e.ErrMsgIterateTagsError, err)
		return false
	}
	return res
}

// SetTag aims to create a given tag in the local repository.
// Args:
// 		repo (*git.Repository): Git repository where the tag must be created.
// 		tag (string): Tag name to be created.
// Returns:
// 		bool: True when tag successfully created, otherwise false.
// 		err: Error whenever unexpected issues happen.
func (s *gitService) SetTag(repo *git.Repository, tag string, authorName string, authorEmail string) (bool, error) {
	tagger := setDefaultSignature(authorName, authorEmail)

	log.Printf("Set tag %s", tag)
	if tagExists(repo, tag) {
		log.Printf(e.ErrMsgTagExists, tag)
		return false, nil
	}

	h, err := s.GetBranchPointedToHead(repo)
	if err != nil {
		log.Fatalf("get HEAD error: %s", err)
		return false, err
	}

	signature := object.Signature{
		Name:  tagger.Name,
		Email: tagger.Email,
		When:  tagger.When,
	}

	log.Printf("Creating tag %s", tag)
	_, err = repo.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger:  &signature,
		Message: fmt.Sprintf("Generated by semantic-release %s", tag),
	})
	if err != nil {
		log.Fatalf("create tag error: %s", err)
		return false, err
	}
	log.Printf("Tag %s successfully created", tag)
	return true, nil
}

// DeleteTag aims to delete a given tag.
// Args:
// 		repo (*git.Repository): Git repository where the tags must be removed from.
// 		tag (string): Tag name to be removed.
// Returns:
// 		bool: True when tag successfully removed, otherwise false.
// 		err: Error whenever unexpected issues happen.
func (s *gitService) DeleteTag(repo *git.Repository, tag string) (bool, error) {
	exists := tagExists(repo, tag)
	fmt.Println(exists)

	if !tagExists(repo, tag) {
		log.Printf(e.ErrMsgTagDoesNotExists, tag)
		return false, errors.New(fmt.Sprintf(e.ErrMsgTagDoesNotExists, tag))
	}
	log.Printf("Deleting tag %s", tag)

	errDelete := repo.DeleteTag(tag)
	if errDelete != nil {
		return false, errDelete
	}
	log.Println("Deleted!")
	return true, nil
}

// PushTags puhs the tags from local to remote repository
// Args:
//		repo (*git.Repository): Git repository where the tags must be pushed to.
// 		auth (*ssh.PublicKeys): SSH key to authenticate while pushing the tags to the remote repository.
// Returns:
// 		error: Returns an error whenever unexpected issues happen.
func (s *gitService) PushTags(repo *git.Repository, auth *ssh.PublicKeys) error {

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       auth,
	}
	err := repo.Push(po)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("origin remote was up to date, no push done")
			return nil
		}
		log.Fatalf("push to remote origin error: %s", err)
		return err
	}
	return nil
}

func New() GitService {
	return &gitService{}
}
