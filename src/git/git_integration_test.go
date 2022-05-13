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
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)
	f.cleanLocalRepo(t)
}

func TestNewGitRepositoryNotFound(t *testing.T) {
	f := setup()
	f.gitLabVersioning.url = "any"
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-set-project")
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : repository not found", err.Error())
}

func TestNewGitEmptyRepositoryError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noBranchProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-branch-project")
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : remote repository is empty", err.Error())
}

func TestNewGitCommitGetChangeHashNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	hash := repo.GetChangeHash()
	tests.AssertEqualValues(t, "a0d3d73a658e905428022c7eca03980569acce5e", hash)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetChangeAuthorNameNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	authorName := repo.GetChangeAuthorName()
	tests.AssertEqualValues(t, "Administrator", authorName)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetChangeAuthorEmailNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	authorEmail := repo.GetChangeAuthorEmail()
	tests.AssertEqualValues(t, "admin@example.com", authorEmail)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetChangeMessageNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	message := repo.GetChangeMessage()
	tests.AssertEqualValues(t, "type: [feat]\r\nmessage: Added requirements.txt file.", message)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetCurrentVersionNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion := repo.GetCurrentVersion()
	tests.AssertEqualValues(t, "1.0.0", currentVersion)
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryNoError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertNoError(t, err)
	f.cleanLocalRepo(t)

	repo, err = f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion = repo.GetCurrentVersion()
	tests.AssertEqualValues(t, newVersion, currentVersion)
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryProtectedTagError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = protectedTagProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "protected-tag-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	err = repo.UpgradeRemoteRepository("0.1.0")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while pushing tag to remote branch due to: command error on refs/tags/1.0.0: pre-receive hook declined", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryAlreadyPushedTagError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertNoError(t, err)

	// push the same tag again
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while pushing tag to remote branch due to: already up-to-date", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryPushToProtectedBranchError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = protectedBranchProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "protected-branch-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while pushing commits to remote repository due to: command error on refs/heads/main: pre-receive hook declined", err.Error())
	f.cleanLocalRepo(t)
}
