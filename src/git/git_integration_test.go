//go:build integration
// +build integration

package git_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/tests"
)

var (
	// These projects have been set previously with a backup which was restored with gitlab after docker compose up.
	// Learn more at `make start-gitlab-env` on the Makefile
	noBranchProject        = fmt.Sprintf("https://%s/dataplatform/no-branch-project.git", host)
	noTagsProject          = fmt.Sprintf("https://%s/dataplatform/no-tags-project.git", host)
	protectedBranchProject = fmt.Sprintf("https://%s/dataplatform/protected-branch-project.git", host)
	protectedTagProject    = fmt.Sprintf("https://%s/dataplatform/protected-tag-project.git", host)
)

func TestNewGitNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)
	f.cleanLocalRepo(t)
}

func TestNewGitRepositoryAlreadyClonedNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)
	tests.AssertNotNil(t, repo)

	_, err = f.newGitService()
	tests.AssertNoError(t, err)
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

	hash := repo.GetChangeHash()
	tests.AssertEqualValues(t, "a0d3d73a658e905428022c7eca03980569acce5e", hash)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetChangeAuthorNameNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	authorName := repo.GetChangeAuthorName()
	tests.AssertEqualValues(t, "Administrator", authorName)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetChangeAuthorEmailNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	authorEmail := repo.GetChangeAuthorEmail()
	tests.AssertEqualValues(t, "admin@example.com", authorEmail)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetChangeMessageNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	message := repo.GetChangeMessage()
	tests.AssertEqualValues(t, "type: [feat]\r\nmessage: Added requirements.txt file.", message)
	f.cleanLocalRepo(t)
}

func TestNewGitCommitGetCurrentVersionNoError(t *testing.T) {
	f := getValidSetup()
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()
	tests.AssertEqualValues(t, "1.0.0", currentVersion)
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryAddChangesError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	f.gitFunctions.errAddToStage = errors.New("no changes to add")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during commit operation due to: no changes to add", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryCommitChangesError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	f.gitFunctions.errCommitChanges = errors.New("commit is old dated")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during commit operation due to: commit is old dated", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryPushError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	f.gitFunctions.errPush = errors.New("nothing to push")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during push operation due to: nothing to push", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositorySetTagError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	f.gitFunctions.errSetTag = errors.New("unable to set new tag")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during set tag operation due to: unable to set new tag", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryPushTagsError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	f.gitFunctions.errPushTag = errors.New("unable to push tags to remote repository")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during push tags operation due to: unable to push tags to remote repository", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryNoError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertNoError(t, err)
	f.cleanLocalRepo(t)

	repo, err = f.newGitService()
	tests.AssertNoError(t, err)

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

	err = repo.UpgradeRemoteRepository("0.1.0")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during push tags operation due to: command error on refs/tags/1.0.0: pre-receive hook declined", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryAlreadyPushedTagError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noTagsProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "no-tags-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertNoError(t, err)

	// push the same tag again
	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during push tags operation due to: already up-to-date", err.Error())
	f.cleanLocalRepo(t)
}

func TestNewGitUpgradeRemoteRepositoryPushToProtectedBranchError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = protectedBranchProject
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "protected-branch-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	err = repo.UpgradeRemoteRepository(newVersion)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error during push operation due to: command error on refs/heads/main: pre-receive hook declined", err.Error())
	f.cleanLocalRepo(t)
}
