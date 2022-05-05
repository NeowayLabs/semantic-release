//go:build unit
// +build unit

package git_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/git"
	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

type GitLabVersioningMock struct {
	url                  string
	destinationDirectory string
	sshKey               string
	replace              bool
}

func printElapsedTimeMock(functionName string) func() {
	return func() {
		fmt.Printf("%s done.", functionName)
	}
}

type fixture struct {
	gitLabVersioning *GitLabVersioningMock
	log              *log.Log
}

func setup() *fixture {
	logger, err := log.New("test", "1.0.0", "debug")
	if err != nil {
		panic(err.Error())
	}
	return &fixture{log: logger, gitLabVersioning: &GitLabVersioningMock{}}
}

func (f *fixture) NewGitService() (git.GitVersioningSystem, error) {
	return git.New(f.log, printElapsedTimeMock, f.gitLabVersioning.url, f.gitLabVersioning.destinationDirectory, f.gitLabVersioning.sshKey, true)
}

func TestNewValidateError(t *testing.T) {
	f := setup()
	gitService, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "url cannot be empty", err.Error())

	f.gitLabVersioning.url = "anyUrl"
	gitService, err = f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "destination directory cannot be empty", err.Error())

	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"

	gitService, err = f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "ssh key cannot be empty", err.Error())
	tests.AssertNil(t, gitService)
}

func TestNewPublicKeyErr(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"
	f.gitLabVersioning.sshKey = "invalidKey"
	gitService, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "ssh: no key found", err.Error())
	tests.AssertNil(t, gitService)
}

func TestNewErrRepositoryNotFound(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "anyUrl"
	f.gitLabVersioning.destinationDirectory = "anyPath"
	f.gitLabVersioning.sshKey = os.Getenv("SSH_INTEGRATION_SEMANTIC")
	gitService, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "repository not found", err.Error())
	tests.AssertNil(t, gitService)
}
