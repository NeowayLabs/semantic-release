//go:build unit
// +build unit

package committype_test

import (
	"testing"

	committype "github.com/NeowayLabs/semantic-release/src/commit-type"
	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

type fixture struct {
	commitType committype.CommitType
}

func setup(t *testing.T) *fixture {
	logger, err := log.New("test", "1.0.0", "debug")
	if err != nil {
		t.Errorf("error while getting log due to %s", err.Error())
	}
	commitType := committype.New(logger)

	return &fixture{commitType: *commitType}
}

func TestGetAll(t *testing.T) {
	f := setup(t)
	expected := []string{"build", "ci", "docs", "fix", "feat", "feature", "feature", "perf", "performance", "refactor", "style", "test", "bc", "breaking", "breaking change", "chore", "skip"}
	actual := f.commitType.GetAll()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetMajorUpgrade(t *testing.T) {
	f := setup(t)
	expected := []string{"bc", "breaking", "breaking change"}
	actual := f.commitType.GetMajorUpgrade()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetMinorUpgrade(t *testing.T) {
	f := setup(t)
	expected := []string{"feat", "feature"}
	actual := f.commitType.GetMinorUpgrade()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetPatchUpgrade(t *testing.T) {
	f := setup(t)
	expected := []string{"build", "ci", "docs", "documentation", "fix", "perf", "performance", "refactor", "style", "test"}
	actual := f.commitType.GetPatchUpgrade()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetSkipVersioning(t *testing.T) {
	f := setup(t)
	expected := []string{"skip", "chore"}
	actual := f.commitType.GetSkipVersioning()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetScopeDefaultSuccess(t *testing.T) {
	f := setup(t)
	actualScope := f.commitType.GetScope("fix: this is the message")
	tests.AssertDeepEqualValues(t, "default", actualScope)
}

func TestGetScopeSuccess(t *testing.T) {
	f := setup(t)
	actualScope := f.commitType.GetScope("fix(scope): this is the message")
	tests.AssertDeepEqualValues(t, "scope", actualScope)
}

func TestGetCommitChangeTypeNotFoundError(t *testing.T) {
	f := setup(t)
	message := "wrong type(scope): This is a sample message"
	actualType, err := f.commitType.GetCommitChangeType(message)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "", actualType)
}

func TestGetCommitChangeTypeSuccess(t *testing.T) {
	f := setup(t)
	expected := "fix"
	message := "fix(scope): This is a sample message"
	actualType, err := f.commitType.GetCommitChangeType(message)
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, expected, actualType)
}

func TestGetCommitChangeTypeNewLinesSuccess(t *testing.T) {
	f := setup(t)
	expected := "feat"
	message := "Merge branch 'sample-branch' into 'master'\n\nfeat(scope): This is a message with new lines.\n\nSee merge request gitgroup/semantic-tests!1"
	actualType, err := f.commitType.GetCommitChangeType(message)
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, expected, actualType)
}
