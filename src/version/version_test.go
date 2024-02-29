//go:build unit
// +build unit

package version_test

import (
	"errors"
	"fmt"
	"testing"

	committype "github.com/NeowayLabs/semantic-release/src/commit-type"
	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/tests"
	"github.com/NeowayLabs/semantic-release/src/version"
)

type fixture struct {
	versionControl *version.VersionControl
}

func setup() *fixture {
	logger, err := log.New("test", "", "info")
	if err != nil {
		errors.New("error while getting new log")
	}

	commitType := committype.New(logger)
	return &fixture{versionControl: version.NewVersionControl(logger, PrintElapsedTimeMock, commitType)}
}

func PrintElapsedTimeMock(what string) func() {
	return func() {
		fmt.Println("print elapsed time")
	}
}

func TestGetNewVersionGetCommitChangeTypeFromMessageError(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("", "")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while finding commit change type within commit message due to: change type not found", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionSplitVersionMajorMinorPatchError(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("feat(scope): this is the message", "1.0.a")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while spliting version into MAJOR.MINOR.PATCH due to: could not convert a to int", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionSplitVersionPathernError(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("feat(scope): this is the message", "1.0")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while spliting version into MAJOR.MINOR.PATCH due to: version must follow the pattern major.minor.patch. I.e.: 1.0.0", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionGetUpgradeTypeError(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("skip(scope): this is the message", "1.0.0")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while getting upgrade type due to: skip is an invalid upgrade change type", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionMajorSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("breaking change(scope): this is the message", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "2.0.0", actualVersion)
}

func TestGetNewVersionMinorSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("feat(scope): this is the message", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.1.0", actualVersion)
}

func TestGetNewVersionPatchSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("fix(scope): this is the message", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.0.1", actualVersion)
}

func TestMustSkipVersioningFalse(t *testing.T) {
	f := setup()
	actualMustSkip := f.versionControl.MustSkipVersioning("fix(scope): this is the message")
	tests.AssertEqualValues(t, false, actualMustSkip)
}

func TestMustSkipVersioningTrue(t *testing.T) {
	f := setup()
	actualMustSkip := f.versionControl.MustSkipVersioning("invalid type(scope): this is the message")
	tests.AssertEqualValues(t, true, actualMustSkip)

	actualMustSkip = f.versionControl.MustSkipVersioning("skip(scope): this is the message")
	tests.AssertEqualValues(t, true, actualMustSkip)
}

func TestGetNewVersionFirstVersionSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("fix(scope): this is the message", "0.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.0.0", actualVersion)

	actualVersion, actualErr = f.versionControl.GetNewVersion("feat(scope): this is the message", "0.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.0.0", actualVersion)
}

func TestGetNewVersionFeatTypeSuccess(t *testing.T) {
	f := setup()
	expected := "1.1.0"
	actualVersion, actualErr := f.versionControl.GetNewVersion("feat: this is the message", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, expected, actualVersion)

	actualVersion, actualErr = f.versionControl.GetNewVersion("feat(default scope): this is the message", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, expected, actualVersion)
}

func TestGetNewVersionAllPatchTypesSuccess(t *testing.T) {
	f := setup()
	patchTypes := []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}
	expected := "1.0.1"

	for _, versionType := range patchTypes {
		actualVersion, actualErr := f.versionControl.GetNewVersion(versionType+": this is the message", "1.0.0")
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, expected, actualVersion)

		actualVersion, actualErr = f.versionControl.GetNewVersion(versionType+"(default scope): this is the message", "1.0.0")
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, expected, actualVersion)
	}
}

func TestGetNewVersionAllMinorTypesSuccess(t *testing.T) {
	f := setup()
	minorTypes := []string{"feat"}
	expected := "1.1.0"

	for _, versionType := range minorTypes {
		actualVersion, actualErr := f.versionControl.GetNewVersion(versionType+": this is the message", "1.0.0")
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, expected, actualVersion)

		actualVersion, actualErr = f.versionControl.GetNewVersion(versionType+"(default scope): this is the message", "1.0.0")
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, expected, actualVersion)
	}
}

func TestGetNewVersionAllMajorTypesSuccess(t *testing.T) {
	f := setup()
	majorTypes := []string{"breaking change", "breaking changes"}
	expected := "2.0.0"

	for _, versionType := range majorTypes {
		actualVersion, actualErr := f.versionControl.GetNewVersion(versionType+": this is the message", "1.0.0")
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, expected, actualVersion)

		actualVersion, actualErr = f.versionControl.GetNewVersion(versionType+"(default scope): this is the message", "1.0.0")
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, expected, actualVersion)
	}
}
