//go:build unit
// +build unit

package git_test

import (
	"testing"

	"github.com/NeowayLabs/semantic-release/src/tests"
)

func TestNewGitEmptyUrlError(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = ""
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "url cannot be empty", err.Error())
}

func TestNewGitEmptyDestinationError(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "any"
	f.gitLabVersioning.destinationDirectory = ""
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "destination directory cannot be empty", err.Error())
}

func TestNewGitEmptyUsernameError(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "any"
	f.gitLabVersioning.destinationDirectory = "any"
	f.gitLabVersioning.username = ""
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "username cannot be empty", err.Error())
}

func TestNewGitEmptyPasswordError(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "any"
	f.gitLabVersioning.destinationDirectory = "any"
	f.gitLabVersioning.username = "any"
	f.gitLabVersioning.password = ""
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "password cannot be empty", err.Error())
}

func TestNewErrRepositoryNotFound(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : repository not found", err.Error())
}
