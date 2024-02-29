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

func TestPrettifyCommitMessageNoMessageError(t *testing.T) {
	f := setup(t)
	message := "feat(scope):"
	prettyMessage, err := f.commitMessageManager.PrettifyCommitMessage(message)
	tests.AssertError(t, err)
	tests.AssertEmpty(t, prettyMessage)
}

func TestPrettifyCommitMessageNewLinesSuccess(t *testing.T) {
	f := setup(t)
	message := "feat(scope): This is a message with new lines."
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
