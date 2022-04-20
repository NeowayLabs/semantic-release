//go:build unit
// +build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	versionService = NewService()
)

func TestConsts(t *testing.T) {
	assert.Equal(t, "MAJOR", major)
	assert.Equal(t, "MINOR", minor)
	assert.Equal(t, "PATCH", patch)
	assert.Equal(t, "%d.%d.%d", versionPattern)
	assert.Equal(t, "__version__", versionPyVariable)
	assert.Equal(t, 150, commitMsgLimitCharacters)
	assert.Equal(t, "@%s", authorChangelogFormat)
}

func TestVars(t *testing.T) {
	// VARS
	assert.Equal(t, []string{"build", "ci", "docs", "fix", "feat", "perf", "refactor", "style", "test", "breaking change", "breaking changes", "skip", "skip versioning", "skip v"}, commitChangeTypes)
	assert.Equal(t, []string{"breaking change", "breaking changes"}, commitChangeTypesMajorUpgrade)
	assert.Equal(t, []string{"feat"}, commitChangeTypesMinorUpgrade)
	assert.Equal(t, []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}, commitChangeTypePatchUpgrade)
	assert.Equal(t, []string{"skip", "skip versioning", "skip v"}, commitTypeSkipVersioning)
}

func TestSplitVersionMajorMinorPatch(t *testing.T) {

	t.Run("No Error", func(t *testing.T) {
		splitedVersion, err := splitVersionMajorMinorPatch("2.1.0")
		assert.NoError(t, err)
		assert.NotNil(t, splitedVersion)
		assert.EqualValues(t, 2, splitedVersion["MAJOR"])
		assert.EqualValues(t, 1, splitedVersion["MINOR"])
		assert.EqualValues(t, 0, splitedVersion["PATCH"])
	})

	t.Run("Error Invalid Version Type", func(t *testing.T) {
		splitedVersion, err := splitVersionMajorMinorPatch("2.1.a")
		assert.Error(t, err)
		assert.EqualValues(t, map[string]int(map[string]int(nil)), splitedVersion)
		assert.EqualValues(t, "could not convert a to int", err.Error())
	})

}

func TestGetUpgradeType(t *testing.T) {

	t.Run("MAJOR Success No Error", func(t *testing.T) {
		upgradeType, err := versionService.GetUpgradeType("breaking change")
		assert.NoError(t, err)
		assert.EqualValues(t, "MAJOR", upgradeType)
	})

	t.Run("MINOR Success No Error", func(t *testing.T) {
		upgradeType, err := versionService.GetUpgradeType("feat")
		assert.NoError(t, err)
		assert.EqualValues(t, "MINOR", upgradeType)
	})

	t.Run("PATCH Success No Error", func(t *testing.T) {
		upgradeType, err := versionService.GetUpgradeType("fix")
		assert.NoError(t, err)
		assert.EqualValues(t, "PATCH", upgradeType)
	})

	t.Run("Invalid Commit Change Type Error", func(t *testing.T) {
		upgradeType, err := versionService.GetUpgradeType("wrong commit type")
		assert.Error(t, err)
		assert.Empty(t, upgradeType)
		assert.EqualValues(t, "wrong commit type is an invalid commit change type", err.Error())
	})
}

func TestGetNewReleaseVersion(t *testing.T) {

	t.Run("Split Version Error", func(t *testing.T) {
		newReleaseVersion, err := versionService.GetNewReleaseVersion("2.1.a", "feat")
		assert.Empty(t, newReleaseVersion)
		assert.Error(t, err)
		assert.EqualValues(t, "error while spliting version into MAJOR.MINOR.PATCH due to: could not convert a to int", err.Error())
	})

	t.Run("Upgrade Type Error", func(t *testing.T) {
		newReleaseVersion, err := versionService.GetNewReleaseVersion("2.1.0", "wrong type")
		assert.Empty(t, newReleaseVersion)
		assert.Error(t, err)
		assert.EqualValues(t, "error while getting upgrade type due to: wrong type is an invalid commit change type", err.Error())
	})

	t.Run("Success MAJOR no Error", func(t *testing.T) {
		newReleaseVersion, err := versionService.GetNewReleaseVersion("2.1.1", "breaking change")
		assert.EqualValues(t, "3.0.0", newReleaseVersion)
		assert.NoError(t, err)
	})

	t.Run("Success MINOR no Error", func(t *testing.T) {
		newReleaseVersion, err := versionService.GetNewReleaseVersion("2.1.1", "feat")
		assert.EqualValues(t, "2.2.0", newReleaseVersion)
		assert.NoError(t, err)
	})

	t.Run("Success PATCH no Error", func(t *testing.T) {
		newReleaseVersion, err := versionService.GetNewReleaseVersion("2.1.1", "fix")
		assert.EqualValues(t, "2.1.2", newReleaseVersion)
		assert.NoError(t, err)
	})
}

func TestPretifyCommitMessage(t *testing.T) {

	t.Run("Short message", func(t *testing.T) {
		expected := "This is the commit message."
		actual, err := versionService.PrettifyCommitMessage("type: [fix]\n, message: This is the commit message.")
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)

		actual, err = versionService.PrettifyCommitMessage(`type: [fix]
		Message: This is the commit message.`)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("Long message", func(t *testing.T) {
		expected := "This is the commit message. it is suppose to have more than consts.commitmsglimitcharacters, and the expected behavior is cutting it, adding three dot..."
		actual, err := versionService.PrettifyCommitMessage("type: [fix], message: This is the commit message. It is suppose to have more than consts.CommitMsgLimitCharacters, and the expected behavior is cutting it, adding three dots at the end of the message.")
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("Error message pattern", func(t *testing.T) {
		actual, err := versionService.PrettifyCommitMessage("type: [fix]\n, This is the commit message.")
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func TestPretifyAuthorEmailForChangelog(t *testing.T) {

	t.Run("No errors", func(t *testing.T) {
		assert.EqualValues(t, "@user.name", versionService.PrettifyAuthorEmailForChangelog("user.name@neoway.com.br"))
	})
}

func TestGetCommitChangeTypeFromMessage(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		changeType, err := versionService.GetCommitChangeTypeFromMessage("type: [feat]\r\nmessage: Added requirements.txt file.")
		assert.NoError(t, err)
		assert.NotNil(t, changeType)
		assert.EqualValues(t, "feat", changeType)
	})

	t.Run("Success", func(t *testing.T) {
		changeType, err := versionService.GetCommitChangeTypeFromMessage("type: [not found]\r\nmessage: Added requirements.txt file.")
		assert.Empty(t, changeType)
		assert.Error(t, err)
		assert.EqualValues(t, "change type not found", err.Error())
	})
}

func TestMustSkipVersioning(t *testing.T) {

	assert.True(t, versionService.MustSkipVersioning("skip"))
	assert.True(t, versionService.MustSkipVersioning("skip v"))
	assert.True(t, versionService.MustSkipVersioning("skip versioning"))
	assert.False(t, versionService.MustSkipVersioning("feat"))

}
