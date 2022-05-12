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
	// These projects have been set previously with a backup which was restored with gitlab after docker compose up.
	// Learn more at `make start-gitlab-env` on the Makefile
	noBranchProject        = "https://gitlab/dataplatform/no-branch-project.git"
	noTagsProject          = "https://gitlab/dataplatform/no-tags-project.git"
	protectedBranchProject = "https://gitlab/dataplatform/protected-branch-project.git"
	protectedTagProject    = "https://gitlab/dataplatform/protected-tag-project.git"
)

func TestNewGitNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)
	f.cleanLocalRepo()
}

func TestNewGitRepositoryNotFound(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "any"
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-set-project")
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	_, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : repository not found", err.Error())
	f.cleanLocalRepo()
}

func TestNewGitEmptyRepositoryError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noBranchProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-branch-project")
	_, err := f.NewGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : remote repository is empty", err.Error())
	f.cleanLocalRepo()
}

func TestNewGitCommitGetChangeHashNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	hash, err := repo.GetChangeHash()
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "a0d3d73a658e905428022c7eca03980569acce5e", hash)
	f.cleanLocalRepo()
}

func TestNewGitCommitGetChangeAuthorNameNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	authorName, err := repo.GetChangeAuthorName()
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "Administrator", authorName)
	f.cleanLocalRepo()
}

func TestNewGitCommitGetChangeAuthorEmailNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	authorEmail, err := repo.GetChangeAuthorEmail()
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "admin@example.com", authorEmail)
	f.cleanLocalRepo()
}

func TestNewGitCommitGetChangeMessageNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	message, err := repo.GetChangeMessage()
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "type: [feat]\r\nmessage: Added requirements.txt file.", message)
	f.cleanLocalRepo()
}

func TestNewGitCommitGetCurrentVersionNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion, err := repo.GetCurrentVersion()
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, "1.0.0", currentVersion)
	f.cleanLocalRepo()
}

func TestNewGitUpgradeRemoteRepositoryNoError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion, err := repo.GetCurrentVersion()
	tests.AssertNoError(t, err)

	newVersion, err := newValidVersion(currentVersion)
	tests.AssertNoError(t, err)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertNoError(t, err)
	f.cleanLocalRepo()

	repo, err = f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion, err = repo.GetCurrentVersion()
	tests.AssertNoError(t, err)
	tests.AssertEqualValues(t, newVersion, currentVersion)
	f.cleanLocalRepo()
}

func TestNewGitUpgradeRemoteRepositoryProtectedTagError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = protectedTagProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "protected-tag-project")
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	err = repo.UpgradeRemoteRepository("0.1.0")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while pushing tag to remote branch due to: command error on refs/tags/1.0.0: pre-receive hook declined", err.Error())
	f.cleanLocalRepo()
}

func TestNewGitUpgradeRemoteRepositoryAlreadyPushedTagError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion, err := repo.GetCurrentVersion()
	tests.AssertNoError(t, err)

	newVersion, err := newValidVersion(currentVersion)
	tests.AssertNoError(t, err)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertNoError(t, err)

	// push the same tag again
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while pushing tag to remote branch due to: already up-to-date", err.Error())
	f.cleanLocalRepo()
}

func TestNewGitUpgradeRemoteRepositoryPushToProtectedBranchError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = protectedBranchProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "protected-branch-project")
	repo, err := f.NewGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion, err := repo.GetCurrentVersion()
	tests.AssertNoError(t, err)

	newVersion, err := newValidVersion(currentVersion)
	tests.AssertNoError(t, err)
	tests.AssertNotEmpty(t, newVersion)

	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while pushing commits to remote repository due to: command error on refs/heads/main: pre-receive hook declined", err.Error())
	f.cleanLocalRepo()
}
