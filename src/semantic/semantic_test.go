//go:build unit
// +build unit

package semantic_test

import (
	"errors"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/semantic"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

type RepositoryVersionControlMock struct {
	hash                 string
	authorName           string
	authorEmail          string
	message              string
	currentVersion       string
	currentChangesInfo   changesInfoMock
	errUpgradeRemoteRepo error
}

func (r *RepositoryVersionControlMock) GetChangeHash() string {
	return r.hash
}
func (r *RepositoryVersionControlMock) GetChangeAuthorName() string {
	return r.authorName
}
func (r *RepositoryVersionControlMock) GetChangeAuthorEmail() string {
	return r.authorEmail
}
func (r *RepositoryVersionControlMock) GetChangeMessage() string {
	return r.message
}
func (r *RepositoryVersionControlMock) GetCurrentVersion() string {
	return r.currentVersion
}

func (r *RepositoryVersionControlMock) UpgradeRemoteRepository(newVersion string) error {
	return r.errUpgradeRemoteRepo
}

type VersionControlMock struct {
	newVersion          string
	errGetNewVersion    error
	mustSkip            bool
	commitChangeType    string
	errCommitChangeType error
}

func (v *VersionControlMock) GetCommitChangeType(commitMessage string) (string, error) {
	return v.commitChangeType, v.errCommitChangeType
}

func (v *VersionControlMock) GetNewVersion(commitMessage string, currentVersion string, upgradeType string) (string, error) {
	return v.newVersion, v.errGetNewVersion
}

func (v *VersionControlMock) MustSkipVersioning(commitMessage string) bool {
	return v.mustSkip
}

type FilesVersionControlMock struct {
	errUpgradeChangeLog       error
	errUpgradeVariableInFiles error
}

func (f *FilesVersionControlMock) UpgradeChangeLog(path, destinationPath string, chageLogInfo interface{}) error {
	return f.errUpgradeChangeLog
}
func (f *FilesVersionControlMock) UpgradeVariableInFiles(filesToUpgrade interface{}, newVersion string) error {
	return f.errUpgradeVariableInFiles
}

type fixture struct {
	rootPath              string
	filesToUpdateVariable interface{}
	repoVersionMock       *RepositoryVersionControlMock
	filesVersionMock      *FilesVersionControlMock
	versionControlMock    *VersionControlMock
}

func setup() *fixture {
	return &fixture{repoVersionMock: &RepositoryVersionControlMock{}, filesVersionMock: &FilesVersionControlMock{}, versionControlMock: &VersionControlMock{}}
}

func (f *fixture) NewSemantic() *semantic.Semantic {
	logger, err := log.New("test", "", "info")
	if err != nil {
		errors.New("error while getting new log")
	}

	version := ""
	return semantic.New(logger, f.rootPath, f.filesToUpdateVariable, f.repoVersionMock, f.filesVersionMock, f.versionControlMock, version)
}

type upgradeFilesMock struct {
	files []upgradeFileMock
}

type upgradeFileMock struct {
	path            string
	destinationPath string
	variableName    string
}

type changesInfoMock struct {
	hash           string
	authorName     string
	authorEmail    string
	message        string
	currentVersion string
	newVersion     string
	changeType     string
}

func (f *fixture) GetValidChangesInfo() changesInfoMock {
	return changesInfoMock{
		hash:           "39a757a0",
		authorName:     "Admin",
		authorEmail:    "admin@admin.com",
		message:        "Any Message",
		currentVersion: "1.0.0",
	}
}

func (f *fixture) GetValidUpgradeFilesInfo() upgradeFilesMock {
	var upgradeFilesInfo upgradeFilesMock
	var upgradeFileInfo upgradeFileMock
	upgradeFileInfo.path = "file/path.py"
	upgradeFileInfo.variableName = "__version__"
	upgradeFileInfo.destinationPath = ""
	upgradeFilesInfo.files = append(upgradeFilesInfo.files, upgradeFileInfo)
	return upgradeFilesInfo
}

func TestGenerateNewReleaseMustSkip(t *testing.T) {
	f := setup()
	f.versionControlMock.mustSkip = true
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertNoError(t, actualErr)
}

func TestGenerateNewReleaseErrorGetNewVersion(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.versionControlMock.errGetNewVersion = errors.New("get new version error")

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while getting new version due to: get new version error", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseErrorUpgradeChangeLog(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.filesVersionMock.errUpgradeChangeLog = errors.New("upgrade changelog error")

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while upgrading changelog file due to: upgrade changelog error", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseErrorUpgradeVariablesInFiles(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.filesVersionMock.errUpgradeVariableInFiles = errors.New("upgrade variables in files error")
	f.filesToUpdateVariable = f.GetValidUpgradeFilesInfo()

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while upgrading variables in files due to: upgrade variables in files error", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseUpgradeRemoteRepositoryError(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.repoVersionMock.errUpgradeRemoteRepo = errors.New("upgrade remote repository error")
	f.filesToUpdateVariable = f.GetValidUpgradeFilesInfo()

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while upgrading remote repository due to: upgrade remote repository error", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseGetCommitChangeError(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.versionControlMock.errCommitChangeType = errors.New("invalid change type")

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()
	tests.AssertError(t, actualErr)
	tests.AssertEqualValues(t, "error while getting commit change type due to: invalid change type", actualErr.Error())
}

func TestGenerateNewReleaseSuccess(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.filesToUpdateVariable = f.GetValidUpgradeFilesInfo()

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()
	tests.AssertNoError(t, actualErr)
}
