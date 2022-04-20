//go:build integration
// +build integration

package git

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/auth"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/stretchr/testify/assert"
)

var (
	gitHost               = "gitlab.integration-tests.com"
	projectName           = "integration-tests"
	projectGroup          = "dataplatform"
	repoUrl               = fmt.Sprintf("git@%s:%s/%s.git", gitHost, projectGroup, projectName)
	repoUrlNotFound       = fmt.Sprintf("git@%s:%s/%s.git", gitHost, projectGroup, "dummyproject")
	authorization         = getPublickKey()
	gitServiceIntegration = New()
)

func getPublickKey() *ssh.PublicKeys {
	keyEnv := os.Getenv("SSH_INTEGRATION_SEMANTIC")
	keyEnv = auth.FormatSSHKey(keyEnv, "#")
	key, errPublicKey := auth.GetPublicKey([]byte(keyEnv))
	if errPublicKey != nil {
		log.Fatalf("error while getting mock ssh public key due to %s", errPublicKey)
	}
	return key
}

func TestCloneRepoToDirectory(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), true, authorization)
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("Error Repo Does Not Exists", func(t *testing.T) {
		repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrlNotFound, fmt.Sprintf("%s/%s", os.Getenv("HOME"), "dummyproject"), true, authorization)
		assert.Error(t, err)
		assert.Nil(t, repo)
	})

	t.Run("Success Repo Already Cloned", func(t *testing.T) {
		repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), false, authorization)
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	})
}

func TestBranchPointedToHead(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), false, authorization)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		branchRef, err := gitServiceIntegration.GetBranchPointedToHead(repo)
		assert.NoError(t, err)
		assert.NotNil(t, branchRef)
		assert.EqualValues(t, "refs/heads/main", branchRef.Name().String())
		assert.EqualValues(t, "a0d3d73a658e905428022c7eca03980569acce5e", branchRef.Hash().String())
		assert.Empty(t, branchRef.Target().String())
		assert.EqualValues(t, "hash-reference", branchRef.Type().String())
	})
}

func TestGetCommitHistory(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), false, authorization)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		branchRef, err := gitServiceIntegration.GetBranchPointedToHead(repo)
		assert.NoError(t, err)

		type Signature struct {
			Name  string
			Email string
		}

		expectedAuthor := Signature{
			Name:  "Administrator",
			Email: "admin@example.com",
		}

		commits, err := gitServiceIntegration.GetCommitHistory(repo, branchRef)
		assert.NoError(t, err)
		assert.NotNil(t, commits)
		assert.EqualValues(t, expectedAuthor.Name, commits[0].Author.Name)
		assert.EqualValues(t, expectedAuthor.Email, commits[0].Author.Email)
		assert.EqualValues(t, "a0d3d73a658e905428022c7eca03980569acce5e", commits[0].Hash.String())
		assert.EqualValues(t, expectedAuthor.Name, commits[0].Committer.Name)
		assert.EqualValues(t, expectedAuthor.Email, commits[0].Committer.Email)
		assert.EqualValues(t, "type: [feat]\r\nmessage: Added requirements.txt file.", commits[0].Message)
	})
}

func TestGetMostRecentCommit(t *testing.T) {

	repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), false, authorization)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	branchRef, err := gitServiceIntegration.GetBranchPointedToHead(repo)
	assert.NoError(t, err)

	commits, err := gitServiceIntegration.GetCommitHistory(repo, branchRef)
	assert.NoError(t, err)
	assert.NotNil(t, commits)

	type Signature struct {
		Name  string
		Email string
	}

	expectedAuthor := Signature{
		Name:  "Administrator",
		Email: "admin@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		var commitList []interface{}
		for _, commit := range commits {
			commitList = append(commitList, commit)
		}

		commit, err := gitServiceIntegration.GetMostRecentCommit(commitList)
		assert.NoError(t, err)
		assert.NotNil(t, commit)
		assert.EqualValues(t, expectedAuthor.Name, commits[0].Author.Name)
		assert.EqualValues(t, expectedAuthor.Email, commits[0].Author.Email)
		assert.EqualValues(t, "a0d3d73a658e905428022c7eca03980569acce5e", commit.Hash.String())
		assert.EqualValues(t, expectedAuthor.Name, commit.Committer.Name)
		assert.EqualValues(t, expectedAuthor.Email, commit.Committer.Email)
		assert.EqualValues(t, "type: [feat]\r\nmessage: Added requirements.txt file.", commit.Message)
	})

	t.Run("Error Parsing Commits Interface", func(t *testing.T) {
		type commitType struct {
			Author   string
			Commiter string
		}

		var commitList []interface{}
		commitList = append(commitList, commitType{Author: "Administrator",
			Commiter: "Administrator",
		})

		commit, err := gitServiceIntegration.GetMostRecentCommit(commitList)
		assert.Nil(t, commit)
		assert.Error(t, err)
		assert.EqualValues(t, "error while casting interface to struct", err.Error())
	})

	t.Run("Error No Commits Found", func(t *testing.T) {
		var commitList []interface{}
		commit, err := gitServiceIntegration.GetMostRecentCommit(commitList)
		assert.Nil(t, commit)
		assert.Error(t, err)
		assert.EqualValues(t, "no commits found", err.Error())
	})
}

func TestGetMostRecentTag(t *testing.T) {

	repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), false, authorization)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	t.Run("Success", func(t *testing.T) {
		recentTag, err := gitServiceIntegration.GetMostRecentTag(repo)
		assert.NoError(t, err)
		assert.NotNil(t, recentTag)
	})
	// TODO: ERROR CASES
}

func TestDeleteTags(t *testing.T) {
	repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), true, authorization)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	t.Run("Tag Does Not Exists", func(t *testing.T) {
		deleted, err := gitServiceIntegration.DeleteTag(repo, "2.0.0")
		assert.False(t, deleted)
		assert.Error(t, err)
		assert.EqualValues(t, "tag 2.0.0 does not exists", err.Error())
	})

	t.Run("Sucess", func(t *testing.T) {
		tagCreated, err := gitServiceIntegration.SetTag(repo, "2.0.0", "Administrator", "admin@example.com")
		assert.NoError(t, err)
		assert.True(t, tagCreated)

		deleted, err := gitServiceIntegration.DeleteTag(repo, "2.0.0")
		assert.True(t, deleted)
		assert.NoError(t, err)

		deleted, err = gitServiceIntegration.DeleteTag(repo, "2.0.0")
		assert.False(t, deleted)
		assert.Error(t, err)
		assert.EqualValues(t, "tag not found", err.Error())
	})

}

func TestSetTag(t *testing.T) {

	repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), true, authorization)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	t.Run("Success", func(t *testing.T) {
		tagCreated, err := gitServiceIntegration.SetTag(repo, "2.0.0", "Administrator", "admin@example.com")
		assert.NoError(t, err)
		assert.True(t, tagCreated)

		deleted, err := gitServiceIntegration.DeleteTag(repo, "2.0.0")
		assert.True(t, deleted)
		assert.NoError(t, err)
	})

}

func TestTagExists(t *testing.T) {
	repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), true, authorization)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	assert.False(t, tagExists(repo, "3.0.0"))

	tagCreated, err := gitServiceIntegration.SetTag(repo, "2.0.0", "Administrator", "admin@example.com")
	assert.NoError(t, err)
	assert.True(t, tagCreated)
	assert.True(t, tagExists(repo, "2.0.0"))
}

func TestPushTag(t *testing.T) {
	repo, err := gitServiceIntegration.CloneRepoToDirectory(repoUrl, fmt.Sprintf("%s/%s", os.Getenv("HOME"), projectName), true, authorization)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	t.Run("Success", func(t *testing.T) {
		tagCreated, err := gitServiceIntegration.SetTag(repo, "1.0.0", "Administrator", "admin@example.com")
		assert.NoError(t, err)
		assert.True(t, tagCreated)

		errPushTags := gitServiceIntegration.PushTags(repo, authorization)
		assert.NoError(t, errPushTags)

		deleted, err := gitServiceIntegration.DeleteTag(repo, "1.0.0")
		assert.True(t, deleted)
		assert.NoError(t, err)
	})
}
