//go:build unit
// +build unit

package commitmessage_test

import (
	"testing"

	commitMessage "github.com/NeowayLabs/semantic-release/src/commit-message"
	committype "github.com/NeowayLabs/semantic-release/src/commit-type"
	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

type fixture struct {
	log                  *log.Log
	commitMessageManager commitMessage.CommitMessage
}

func setup(t *testing.T) *fixture {
	logger, err := log.New("test", "1.0.0", "debug")
	if err != nil {
		t.Errorf("error while getting log due to %s", err.Error())
	}

	commitType := committype.New(logger)
	commitMessageMenager := commitMessage.New(logger, commitType)

	return &fixture{log: logger, commitMessageManager: *commitMessageMenager}
}

func TestPrettifyCommitMessageNoMessageEmptyError(t *testing.T) {
	f := setup(t)
	message := "feat(scope):"
	prettyMessage, err := f.commitMessageManager.PrettifyCommitMessage(message)
	tests.AssertError(t, err)
	tests.AssertEmpty(t, prettyMessage)
}

func TestPrettifyCommitMessageNewLinesSuccess(t *testing.T) {
	f := setup(t)
	message := "Merge branch 'sample-branch' into 'master'\n\nfeat(scope): This is a message with new lines.\n\nSee merge request gitgroup/semantic-tests!1"
	prettyMessage, err := f.commitMessageManager.PrettifyCommitMessage(message)
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "This is a message with new lines.", prettyMessage)
}

func TestPrettifyCommitMessageCutSuccess(t *testing.T) {
	f := setup(t)
	message := "feat: This is a long message to write to CHANGELOG.md file. Bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo cut here."
	prettyMessage, err := f.commitMessageManager.PrettifyCommitMessage(message)
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "This is a long message to write to CHANGELOG.md file. Bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo ...", prettyMessage)
}

func TestIsValidMessageSuccess(t *testing.T) {
	f := setup(t)
	message := "Merge branch 'sample-branch' into 'master'\n\nfeat(scope): This is a message with new lines.\n\nSee merge request gitgroup/semantic-tests!1"
	actual := f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "feat(scope): This is a message"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "feat: This is a message"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)
}

func TestIsValidMessageFalse(t *testing.T) {
	f := setup(t)
	message := "Merge branch 'sample-branch' into 'master'\n\nfeat(scope) This is a message with new lines.\n\nSee merge request gitgroup/semantic-tests!1"
	actual := f.commitMessageManager.IsValidMessage(message)
	tests.AssertFalse(t, actual)

	message = "feat(scope):"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertFalse(t, actual)

	message = "feat:"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertFalse(t, actual)

	message = "feat This is a message with new lines"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertFalse(t, actual)

	message = ""
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertFalse(t, actual)

	message = "wrong type(scope): This is a message"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertFalse(t, actual)
}

func TestIsValidMessageMergeMasterBranchSuccess(t *testing.T) {
	f := setup(t)
	message := "first message row \n Merge remote-tracking branch 'origin/master' into something \n last message row"
	actual := f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "first message row \n Merge branch 'master' into something \n last message row"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "first message row \n Merge branch 'master' of something \n last message row"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "first message row \n Merge remote-tracking branch 'origin/main' into something \n last message row"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "first message row \n Merge branch 'main' into something \n last message row"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)

	message = "first message row \n Merge branch 'main' of something \n last message row"
	actual = f.commitMessageManager.IsValidMessage(message)
	tests.AssertTrue(t, actual)
}
