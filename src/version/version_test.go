//go:build unit
// +build unit

package version_test

import (
	"errors"
	"fmt"
	"testing"

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

	return &fixture{versionControl: version.NewVersionControl(logger, PrintElapsedTimeMock)}
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
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[feat]", "1.0.a")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while spliting version into MAJOR.MINOR.PATCH due to: could not convert a to int", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionSplitVersionPathernError(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[feat]", "1.0")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while spliting version into MAJOR.MINOR.PATCH due to: version must follow the pattern major.minor.patch. I.e.: 1.0.0", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionGetUpgradeTypeError(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[skip]", "1.0.0")
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while getting upgrade type due to: skip is an invalid upgrade change type", actualErr.Error())
	tests.AssertEmpty(t, actualVersion)
}

func TestGetNewVersionMajorSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[breaking change]", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "2.0.0", actualVersion)
}

func TestGetNewVersionMinorSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[feat]", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.1.0", actualVersion)
}

func TestGetNewVersionPatchSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[fix]", "1.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.0.1", actualVersion)
}

func TestGetNewVersionPatchSuccessLot(t *testing.T) {
	f := setup()
	testcases := map[string]string{
		"type:[fix]": "1.0.0",
		"type:[fix]": "1.0.0",
	}
	expected := []string{
		"1.0.1",
	}
	for k, v := range testcases {
		actualVersion, actualErr := f.versionControl.GetNewVersion(k, v)
		tests.AssertNoError(t, actualErr)
		tests.AssertEqualValues(t, "1.0.1", actualVersion)
	}

}

func TestMustSkipVersioningFalse(t *testing.T) {
	f := setup()
	actualMustSkip := f.versionControl.MustSkipVersioning("type: [fix]")
	tests.AssertEqualValues(t, false, actualMustSkip)
}

func TestMustSkipVersioningTrue(t *testing.T) {
	f := setup()
	actualMustSkip := f.versionControl.MustSkipVersioning("type: [anything]")
	tests.AssertEqualValues(t, true, actualMustSkip)

	actualMustSkip = f.versionControl.MustSkipVersioning("type: [skip]")
	tests.AssertEqualValues(t, true, actualMustSkip)
}

func TestGetNewVersionFirstVersionSuccess(t *testing.T) {
	f := setup()
	actualVersion, actualErr := f.versionControl.GetNewVersion("type:[fix]", "0.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.0.0", actualVersion)

	actualVersion, actualErr = f.versionControl.GetNewVersion("type:[feat]", "0.0.0")
	tests.AssertNoError(t, actualErr)
	tests.AssertEqualValues(t, "1.0.0", actualVersion)
}
