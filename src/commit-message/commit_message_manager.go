package commitmessage

import (
	"errors"
	"fmt"
	"strings"
)

const (
	messageTag = "message:"
)

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
	Warn(s string, args ...interface{})
}

type CommitType interface {
	GetAll() []string
	GetMajorUpgrade() []string
	GetMinorUpgrade() []string
	GetPatchUpgrade() []string
	GetSkipVersioning() []string
	GetCommitChangeType(commitMessage string) (string, error)
}

type CommitMessage struct {
	log        Logger
	commitType CommitType
}

func (f *CommitMessage) findMessageTag(commitMessage string) bool {
	return strings.Contains(strings.ToLower(commitMessage), messageTag)
}

func (f *CommitMessage) isMessageLongerThanLimit(message string) bool {
	return len(message) >= 150
}

func (f *CommitMessage) upperFirstLetterOfSentence(text string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(text[:1]), text[1:])
}

// prettifyCommitMessage aims to keep a short message based on the commit message, removing extra information such as commit type.
// Args:
//
//	commitMessage (string): Full commit message.
//
// Returns:
//
//	string: Returns a commit message with limmited number of characters.
//	err: Error whenever unexpected issues happen.
func (f *CommitMessage) PrettifyCommitMessage(commitMessage string) (string, error) {
	splitedMessage := strings.Split(commitMessage, "\n")

	message := ""
	for _, row := range splitedMessage {
		if row == "" {
			continue
		}
		index := strings.Index(row, ":")
		commitTypeScope := strings.ToLower(row[:index])

		for _, changeType := range f.commitType.GetAll() {
			if strings.Contains(commitTypeScope, changeType) {
				index := strings.Index(row, ":")
				message = strings.TrimSpace(strings.Replace(row[index:], ":", "", 1))
			}
		}
	}

	if message == "" {
		return "", errors.New("commit message is empty")
	}

	if f.isMessageLongerThanLimit(message) {
		message = fmt.Sprintf("%s...", message[:150])
	}

	return f.upperFirstLetterOfSentence(message), nil
}

func (f *CommitMessage) IsValidMessage(message string) bool {
	index := strings.Index(message, ":")

	if index == -1 {
		f.log.Error("commit message out of pattern")
		return false
	}

	if message == "" || message[index:] == "" {
		f.log.Error("commit message cannot be empty")
		return false
	}

	_, err := f.commitType.GetCommitChangeType(message)
	if err != nil {
		if err.Error() == "change type not found" {
			f.log.Error("change type not found")
		}
		return false
	}

	return true
}

func New(log Logger, commitType CommitType) *CommitMessage {
	return &CommitMessage{
		log:        log,
		commitType: commitType,
	}
}
