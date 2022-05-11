//go:build unit
// +build unit

package git_test

import (
	"testing"

	"github.com/NeowayLabs/semantic-release/src/tests"
)

func TestNewValidateError(t *testing.T) {
	f := setup()
	_, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "url cannot be empty", err.Error())

	f.gitLabVersioning.url = "anyUrl"
	_, err = f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "destination directory cannot be empty", err.Error())

	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"
	_, err = f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "username cannot be empty", err.Error())

	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"
	f.gitLabVersioning.username = "root"
	_, err = f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "password cannot be empty", err.Error())
}

func TestNewErrRepositoryNotFound(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	_, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : repository not found", err.Error())
}
