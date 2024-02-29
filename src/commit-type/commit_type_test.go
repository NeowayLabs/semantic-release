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
	expected := []string{"build", "ci", "docs", "fix", "feat", "perf", "refactor", "style", "test", "breaking change", "breaking changes", "skip"}
	actual := f.commitType.GetAll()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetMajorUpgrade(t *testing.T) {
	f := setup(t)
	expected := []string{"breaking change", "breaking changes"}
	actual := f.commitType.GetMajorUpgrade()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetMinorUpgrade(t *testing.T) {
	f := setup(t)
	expected := []string{"feat"}
	actual := f.commitType.GetMinorUpgrade()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetPatchUpgrade(t *testing.T) {
	f := setup(t)
	expected := []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}
	actual := f.commitType.GetPatchUpgrade()

	tests.AssertDeepEqualValues(t, expected, actual)
}

func TestGetSkipVersioning(t *testing.T) {
	f := setup(t)
	expected := []string{"skip"}
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