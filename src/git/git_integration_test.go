//go:build integration
// +build integration

package git_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/tests"
)

func TestNewErrRepositoryNoError(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "https://gitlab.integration-tests.com/dataplatform/integration-tests.git"
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "semantic-release")
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)
}
