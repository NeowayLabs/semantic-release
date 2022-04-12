package errorsutils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstantsErrorsMessages(t *testing.T) {
	assert.EqualValues(t, "error while clonning repository", ErrMsgCloneFail)
	assert.EqualValues(t, "error while retrieving the branch pointed by HEAD", ErrMsgRetrievingBranchHead)
	assert.EqualValues(t, "error while retrieving the commit history", ErrMsgRetrievingCommitHistory)
	assert.EqualValues(t, "error while iterating over the commits", ErrMsgIteratingCommitHistory)
	assert.EqualValues(t, "change type not found", ErrMsgGetChangeType)
	assert.EqualValues(t, "%s is an invalid commit change type", ErrMsgInvalidCommitChangeType)
	assert.EqualValues(t, "error while retrieving tags from repository", ErrMsgRetrievingTags)
	assert.EqualValues(t, "error while getting most recent tag", ErrMsgGettingMostRecentTag)
	assert.EqualValues(t, "error while finding commit change type within commit message", ErrMsgGetCommitChangeTypeFromMessage)
	assert.EqualValues(t, "error while getting most recent commit", ErrMsgGettingMostRecentCommit)
	assert.EqualValues(t, "no commits found", ErrMsgNoCommitsFound)
	assert.EqualValues(t, "repository was already cloned", ErrMsgRepoAlreadyCloned)
	assert.EqualValues(t, "%s due to: %s", ErrMsgDueTo)
	assert.EqualValues(t, "error while upgrading setup.py file due to: %s", ErrMsgUpgradeSetupPython)
	assert.EqualValues(t, "error while upgrading CHANGELOG.md file due to: %s", ErrMsgUpgradeChangelog)
	assert.EqualValues(t, "no such file or directory", ErrMsgNoSuchFileOrDirectory)
	assert.EqualValues(t, "\n\nerror while writing %s file with new version %s due to: %s", ErrMsgWritingUpgradedFile)
	assert.EqualValues(t, "\n\nerror while scanning file: %s due to: %s", ErrMsgScanningFile)
	assert.EqualValues(t, "commit message has no tag 'message:'", ErrMsgGetCommitMsgPatternNotFound)
	assert.EqualValues(t, "error while getting new release version", ErrMsgNewReleaseVersion)

}

func TestVerifyErrors(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	err := errors.New("internal server error")
	fmt.Print(err)
	// VerifyError(err, "something went wrong")
	fmt.Println(string(buf.Bytes()))
	t.Log(buf.String())
}

func TestError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		err := errors.New("error cause")
		assert.True(t, Error(err, "error while doing something"))
		assert.EqualValues(t, "error cause", err.Error())

	})

	t.Run("No Error", func(t *testing.T) {
		assert.False(t, Error(nil, "error while doing something"))
	})
}

func TestHasErrorPanic(t *testing.T) {
	err := errors.New("error cause")
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("HasError should have panicked!")
			}
		}()
		// This function should cause a panic
		HasError(err, "error while doing something")
	}()

}
