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
	// Learn more at `make start-env` on the Makefile
	noBranchProject         = fmt.Sprintf("https://%s/dataplatform/no-branch-project.git", host)
	noTagsProject           = fmt.Sprintf("https://%s/dataplatform/no-tags-project.git", host)
	protectedBranchProject  = fmt.Sprintf("https://%s/dataplatform/protected-branch-project.git", host)
	protectedTagProject     = fmt.Sprintf("https://%s/dataplatform/protected-tag-project.git", host)
	alphaTagProject         = fmt.Sprintf("https://%s/dataplatform/alpha-tag-project.git", host)
	alphaNumericTagProject  = fmt.Sprintf("https://%s/dataplatform/alpha-numeric-tag-project.git", host)
	tagsOutOfPatternProject = fmt.Sprintf("https://%s/dataplatform/tags-out-of-pattern-project.git", host)
	greatNumbersTagProject  = fmt.Sprintf("https://%s/dataplatform/great-numbers-tag-project.git", host)
	disorderedTagsProject   = fmt.Sprintf("https://%s/dataplatform/disordered-tags-prject.git", host)
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-set-project")
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	_, err := f.newGitService()
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while initiating git package due to : repository not found", err.Error())
}

func TestNewGitEmptyRepositoryError(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = noBranchProject
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-branch-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	branchHead := repo.BranchHead()

	currentVersion := repo.GetCurrentVersion()

	newVersion := newValidVersion(t, currentVersion)
	tests.AssertNotEmpty(t, newVersion)

	// push the tag once
	err = repo.UpgradeRemoteRepository(newVersion)
	newBranchHead := repo.BranchHead()
	tests.AssertNoError(t, err)
	tests.AssertDiffValues(t, branchHead, newBranchHead)
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("protected-tag-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("no-tags-project")
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
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("protected-branch-project")
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

func TestNewGitGetCurrentVersionFromRepoWithAlphaCharacteres(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = alphaTagProject
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("alpha-tag-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()
	tests.AssertEqualValues(t, "0.0.0", currentVersion)
	f.cleanLocalRepo(t)
}

func TestNewGitGetCurrentVersionFromRepoWithAlphaAndNumericCharacteres(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = alphaNumericTagProject
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("alpha-numeric-tag-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	currentVersion := repo.GetCurrentVersion()
	tests.AssertEqualValues(t, "1.0.1", currentVersion)
	f.cleanLocalRepo(t)
}

func TestNewGitGetCurrentVersionFromRepoWithTagsOutOfPattern(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = tagsOutOfPatternProject
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("tags-out-of-pattern-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	result, err := repo.GetMostRecentTag()
	tests.AssertNil(t, err)
	tests.AssertEqualValues(t, "2.1.0", result)
	f.cleanLocalRepo(t)
}

func TestNewGitGetCurrentVersionFromRepoWithGreatNumbersTag(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = greatNumbersTagProject
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("great-numbers-tag-project")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	result, err := repo.GetMostRecentTag()
	tests.AssertNil(t, err)
	tests.AssertEqualValues(t, "3.0.2", result)
	f.cleanLocalRepo(t)
}

func TestNewGitGetCurrentVersionFromRepoWithDisorderedTags(t *testing.T) {
	f := getValidSetup()
	f.gitLabVersioning.url = disorderedTagsProject
	f.gitLabVersioning.destinationDirectory = getDestinationDirectory("disordered-tags-prject")
	repo, err := f.newGitService()
	tests.AssertNoError(t, err)

	result, err := repo.GetMostRecentTag()
	tests.AssertNil(t, err)
	tests.AssertEqualValues(t, "20.1.0", result)
	f.cleanLocalRepo(t)
}

func getDestinationDirectory(repo string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("HOME"), repo)
}

// func TestGetChangelogWorks(t *testing.T) {
// 	f := getValidSetup()
// 	defer f.cleanLocalRepo(t)
// 	f.gitLabVersioning.url = greatNumbersTagProject
// 	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "get-changelog")
// 	repo, err := f.newGitService()
// 	file, err := os.Create(fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), "get-changelog", "CHANGELOG"))
// 	if err != nil {
// 		fmt.Println("Error to create file", err)
// 		return
// 	}
// 	defer file.Close()
// 	tests.AssertNoError(t, err)
// 	fmt.Fprintf(file, "type [fix], message: testing message")

// 	result, err := repo.GetChangelog()
// 	tests.AssertNil(t, err)
// 	tests.AssertNotNil(t, result)
// }

// func TestGetChangelogChangesWorks(t *testing.T) {
// 	f := getValidSetup()
// 	f.gitLabVersioning.url = greatNumbersTagProject
// 	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "update-changelog")
// 	repo, err := f.newGitService()
// 	file, err := os.Create(fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), "update-changelog", "CHANGELOG"))
// 	if err != nil {
// 		fmt.Println("Error to create file", err)
// 		return
// 	}
// 	defer file.Close()
// 	tests.AssertNoError(t, err)
// 	changelogSample := getChangelogSample()
// 	fmt.Fprintf(file, changelogSample)

// 	result, err := repo.GetChangelogChanges()
// 	tests.AssertNil(t, err)
// 	tests.AssertNotNil(t, result)
// 	tests.AssertEqualValues(t, "type: [fix], message: testing message", result)
// 	f.cleanLocalRepo(t)
// }

// func TestGetChangelogChangesError(t *testing.T) {
// 	f := getValidSetup()
// 	f.gitLabVersioning.url = greatNumbersTagProject
// 	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "update-changelog")
// 	repo, err := f.newGitService()
// 	file, err := os.Create(fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), "update-changelog", "CHANGELOG"))
// 	if err != nil {
// 		fmt.Println("Error to create file", err)
// 		return
// 	}
// 	defer file.Close()
// 	tests.AssertNoError(t, err)
// 	changelogSample := getChangelogSampleWithTwoChanges()
// 	fmt.Fprintf(file, changelogSample)

// 	result, err := repo.GetChangelogChanges()
// 	tests.AssertEmpty(t, result)
// 	tests.AssertNotNil(t, err)
// 	f.cleanLocalRepo(t)
// }

// func getChangelogSample() string {
// 	return `type: [fix], message: testing message

// 	## v1.9.9
// 	type: [fix], message: v199`
// }

// func getChangelogSampleWithTwoChanges() string {
// 	return `type: [fix], message: testing message
// 	type: [feat], message: feat message

// 	## v1.9.9
// 	type: [fix], message: v199`
// }
