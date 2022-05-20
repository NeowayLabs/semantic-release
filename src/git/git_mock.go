package git

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Git interface {
	GetBranchPointedToHead() (*plumbing.Reference, error)
	GetCommitHistory() ([]*object.Commit, error)
	GetMostRecentCommit() (CommitInfo, error)
	GetAllTags() ([]object.Tag, error)
	GetMostRecentTag() (string, error)
	AddToStage() error
	CommitChanges(newReleaseVersion string) error
	Push() error
	TagExists(tag string) (bool, error)
	SetTag(tag string) error
	PushTags() error
}

// substituteFunctions aims to replace the package functions with the injected ones.
// It can be used to pass mock funcions or specific implementations for a given function.
func (g *GitVersioning) substituteFunctions(newGit Git) {

	branchHead, err := newGit.GetBranchPointedToHead()
	if err != nil || branchHead != nil {
		g.git.getBranchPointedToHead = newGit.GetBranchPointedToHead
	}

	commitHistory, err := newGit.GetCommitHistory()
	if err != nil || commitHistory != nil {
		g.git.getCommitHistory = newGit.GetCommitHistory
	}

	var emptyCommitInfo CommitInfo
	mostRecentCommit, err := newGit.GetMostRecentCommit()
	if err != nil || mostRecentCommit != emptyCommitInfo {
		g.git.getMostRecentCommit = newGit.GetMostRecentCommit
	}

	allTags, err := newGit.GetAllTags()
	if err != nil || allTags != nil {
		g.git.getAllTags = newGit.GetAllTags
	}

	mostRecentTag, err := newGit.GetMostRecentTag()
	if err != nil || mostRecentTag != "" {
		g.git.getMostRecentTag = newGit.GetMostRecentTag
	}

	if err := newGit.AddToStage(); err != nil {
		g.git.addToStage = newGit.AddToStage
	}

	if err := newGit.CommitChanges(""); err != nil {
		g.git.commitChanges = newGit.CommitChanges
	}

	if err := newGit.Push(); err != nil {
		g.git.push = newGit.Push
	}

	_, err = newGit.TagExists("")
	if err != nil {
		g.git.tagExists = newGit.TagExists
	}

	if err := newGit.SetTag(""); err != nil {
		g.git.setTag = newGit.SetTag
	}

	if err := newGit.PushTags(); err != nil {
		g.git.pushTags = newGit.PushTags
	}
}

func NewMock(log Logger, printElapsedTime ElapsedTime, url, username, password, destinationDirectory string, git Git) (*GitVersioning, error) {

	gitLabVersioning, err := New(log, printElapsedTime, url, username, password, destinationDirectory)
	if err != nil {
		return nil, err
	}

	// mock or new implementation
	gitLabVersioning.substituteFunctions(git)

	if err := gitLabVersioning.initialize(); err != nil {
		return nil, err
	}

	return gitLabVersioning, nil
}
