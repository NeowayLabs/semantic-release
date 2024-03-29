package committype

import (
	"errors"
	"regexp"
	"strings"
)

type Logger interface {
	Info(s string, args ...interface{})
}

type CommitType struct {
	log Logger
}

func (c *CommitType) GetAll() []string {
	return []string{"build", "ci", "docs", "fix", "feat", "feature", "feature", "perf", "performance", "refactor", "style", "test", "bc", "breaking", "breaking change", "chore", "skip"}
}

func (c *CommitType) GetMajorUpgrade() []string {
	return []string{"bc", "breaking", "breaking change"}
}

func (c *CommitType) GetMinorUpgrade() []string {
	return []string{"feat", "feature"}
}

func (c *CommitType) GetPatchUpgrade() []string {
	return []string{"build", "ci", "docs", "documentation", "fix", "perf", "performance", "refactor", "style", "test"}
}

func (c *CommitType) GetSkipVersioning() []string {
	return []string{"skip", "chore"}
}

func (c *CommitType) isValidCommitType(commitTypeScope string) bool {
	for _, changeType := range c.GetAll() {
		if strings.Contains(commitTypeScope, changeType) {
			return true
		}
	}
	return false
}

// GetScope get the commit scope from Message
// I.e.:
//
//	fix(any): Commit subject here.
//
// Output: any
func (c *CommitType) GetScope(commitMessage string) string {
	c.log.Info("getting commit scope from message %s", commitMessage)
	splitedMessage := strings.Split(commitMessage, "\n")
	re := regexp.MustCompile(`\((.*?)\)`)

	for _, row := range splitedMessage {
		index := strings.Index(row, ":")

		if c.IndexNotFound(index) || row == " " {
			continue
		}
		commitTypeScope := strings.ToLower(row[:index])

		if c.isValidCommitType(commitTypeScope) {
			found := re.FindAllString(row, -1)
			for _, element := range found {
				element = strings.Trim(element, "(")
				element = strings.Trim(element, ")")
				return element
			}
		}
	}

	return "default"
}

func (c *CommitType) IndexNotFound(index int) bool {
	return index == -1
}

// GetCommitChangeType get the commit type from Message
// I.e.:
//
//	fix(scope?): Commit subject here.
//
// Output: fix
func (c *CommitType) GetCommitChangeType(commitMessage string) (string, error) {
	c.log.Info("getting commit type from message %s", commitMessage)
	splitedMessage := strings.Split(commitMessage, "\n")

	for _, row := range splitedMessage {
		index := strings.Index(row, ":")

		if c.IndexNotFound(index) || row == "" {
			continue
		}

		commitTypeScope := strings.ToLower(row[:index])

		for _, changeType := range c.GetAll() {
			if strings.Contains(commitTypeScope, changeType) {
				return changeType, nil
			}
		}
	}

	return "", errors.New("change type not found")
}

func New(log Logger) *CommitType {
	return &CommitType{
		log: log,
	}
}
