//go:build integration
// +build integration

package git_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/tests"
)

var (
	validUrl = "https://gitlab/dataplatform/integration-tests.git"
)

func TestNewRepositoryNoError(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = validUrl
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "integration-tests")
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)
}
