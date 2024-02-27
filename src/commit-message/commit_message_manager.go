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

type CommitMessage struct {
	log Logger
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

func (f *CommitMessage) getMessage(messageRow string) (string, error) {
	startPosition := strings.Index(messageRow, messageTag) + len(messageTag)

	if startPosition-1 == len(messageRow)-1 {
		return "", errors.New("message not found")
	}

	message := strings.TrimSpace(messageRow[startPosition:])
	if strings.ReplaceAll(message, " ", "") == "" {
		return "", errors.New("message not found")
	}

	return message, nil
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

	var message string
	splitedMessage := strings.Split(commitMessage, "\n")

	for _, row := range splitedMessage {
		row := strings.ToLower(row)
		if f.findMessageTag(row) {

			currentMessage, err := f.getMessage(row)
			if err != nil {
				return "", fmt.Errorf("error while getting message due to: %w", err)
			}
			message = currentMessage
		}
	}

	if message == "" {
		return "", errors.New("commit message has no tag 'message:'")
	}

	if f.isMessageLongerThanLimit(message) {
		message = fmt.Sprintf("%s...", message[:150])
	}

	return f.upperFirstLetterOfSentence(message), nil
}

func New(log Logger) *CommitMessage {
	return &CommitMessage{
		log: log,
	}
}
