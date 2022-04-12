package consts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalConstsAndVars(t *testing.T) {
	// CONSTS
	assert.Equal(t, "MAJOR", Major)
	assert.Equal(t, "MINOR", Minor)
	assert.Equal(t, "PATCH", Patch)
	assert.Equal(t, "%d.%d.%d", VersionPattern)
	assert.Equal(t, "__version__", VersionPyVariable)
	assert.Equal(t, "CHANGELOG.md", ChangeLogDefaultFile)
	assert.Equal(t, "setup.py", SetupPythonDefaultFile)
	assert.Equal(t, 150, CommitMsgLimitCharacters)
	assert.Equal(t, "@%s", AuthorChangelogFormat)
	assert.Equal(t, "[%s](https://%s/%s/%s/-/commit/%s)", ChangeLogCommitHashLinkFormat)

	// VARS
	assert.Equal(t, []string{"build", "ci", "docs", "fix", "feat", "perf", "refactor", "style", "test", "breaking change", "skip", "skip versioning", "skip v"}, CommitChangeTypes)
	assert.Equal(t, []string{"breaking change"}, CommitChangeTypesMajorUpgrade)
	assert.Equal(t, []string{"feat"}, CommitChangeTypesMinorUpgrade)
	assert.Equal(t, []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}, CommitChangeTypePatchUpgrade)
	assert.Equal(t, []string{"skip", "skip versioning", "skip v"}, CommitTypeSkipVersioning)
}
