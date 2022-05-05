//go:build unit
// +build unit

package semantic_test

import (
	"errors"
	"os"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/semantic"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

const (
	shouldbeEqualMessage = "Values should be equal:\n Actual: %v\nExpected: %v"
)

type RepositoryVersionControlMock struct {
	hash       string
	errGetHash error

	authorName    string
	errAuthorName error

	authorEmail    string
	errAuthorEmail error

	message       string
	errGetMessage error

	currentVersion       string
	errGetCurrentVersion error

	currentChangesInfo    changesInfoMock
	errCurrentChangesInfo error

	errUpgradeRemoteRepo error
}

func (r *RepositoryVersionControlMock) GetChangeHash() (string, error) {
	return r.hash, r.errGetHash
}
func (r *RepositoryVersionControlMock) GetChangeAuthorName() (string, error) {
	return r.authorName, r.errAuthorName
}
func (r *RepositoryVersionControlMock) GetChangeAuthorEmail() (string, error) {
	return r.authorEmail, r.errAuthorEmail
}
func (r *RepositoryVersionControlMock) GetChangeMessage() (string, error) {
	return r.message, r.errGetMessage
}
func (r *RepositoryVersionControlMock) GetCurrentVersion() (string, error) {
	return r.currentVersion, r.errGetCurrentVersion
}

func (r *RepositoryVersionControlMock) getChangesInformation() (*changesInfoMock, error) {
	return &r.currentChangesInfo, r.errCurrentChangesInfo
}

func (r *RepositoryVersionControlMock) UpgradeRemoteRepository(newVersion string) error {
	return r.errUpgradeRemoteRepo
}

type VersionControlMock struct {
	newVersion       string
	errGetNewVersion error
	mustSkip         bool
}

func (v *VersionControlMock) GetNewVersion(commitMessage string, currentVersion string) (string, error) {
	return v.newVersion, v.errGetNewVersion
}

func (v *VersionControlMock) MustSkipVersioning(commitMessage string) bool {
	return v.mustSkip
}

type FilesVersionControlMock struct {
	errUpgradeChangeLog       error
	errUpgradeVariableInFiles error
}

func (f *FilesVersionControlMock) UpgradeChangeLog(path string, chagelogInfo interface{}, newVersion string) error {
	return f.errUpgradeChangeLog
}
func (f *FilesVersionControlMock) UpgradeVariableInFiles(filesInfo interface{}, newVersion string) error {
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
	return semantic.New(f.rootPath, f.filesToUpdateVariable, f.repoVersionMock, f.filesVersionMock, f.versionControlMock)
}

type upgradeFilesMock struct {
	files []upgradeFileMock
}

type upgradeFileMock struct {
	path         string
	variableName string
}

type changesInfoMock struct {
	hash           string
	authorName     string
	authorEmail    string
	message        string
	currentVersion string
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
	upgradeFilesInfo.files = append(upgradeFilesInfo.files, upgradeFileInfo)
	return upgradeFilesInfo
}

func TestGenerateNewReleaseErrorGetChangeHash(t *testing.T) {
	f := setup()
	f.repoVersionMock.errGetHash = errors.New("test")
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while getting changes information due to: error getting hash: test", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseErrorGetChangeAuthorName(t *testing.T) {
	f := setup()
	f.repoVersionMock.errAuthorName = errors.New("test")
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while getting changes information due to: error getting author name: test", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseErrorGetChangeAuthorEmail(t *testing.T) {
	f := setup()
	f.repoVersionMock.errAuthorEmail = errors.New("test")
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while getting changes information due to: error getting author email: test", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseErrorGetChangeMessage(t *testing.T) {
	f := setup()
	f.repoVersionMock.errGetMessage = errors.New("test")
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while getting changes information due to: error getting message: test", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseMustSkip(t *testing.T) {
	f := setup()
	f.versionControlMock.mustSkip = true
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertNoError(t, actualErr)
}

func TestGenerateNewReleaseErrorGetCurrentVersion(t *testing.T) {
	f := setup()
	f.repoVersionMock.errGetCurrentVersion = errors.New("test")
	f.rootPath = os.Getenv("HOME")
	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error while getting changes information due to: error getting current version: test", actualErr.Error())
	tests.AssertError(t, actualErr)
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

func TestGenerateNewReleaseErrorUpgradeVariablesInFilesMarshalError(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.filesToUpdateVariable = make(chan int)

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error marshalling files to uptade information", actualErr.Error())
	tests.AssertError(t, actualErr)
}

func TestGenerateNewReleaseErrorUpgradeVariablesInFilesUnmarshalError(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.filesToUpdateVariable = "anything"

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()

	tests.AssertEqualValues(t, "error unmarshalling files to uptade information", actualErr.Error())
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

func TestGenerateNewReleaseSuccess(t *testing.T) {
	f := setup()
	f.repoVersionMock.currentChangesInfo = f.GetValidChangesInfo()
	f.filesToUpdateVariable = f.GetValidUpgradeFilesInfo()

	semanticService := f.NewSemantic()
	actualErr := semanticService.GenerateNewRelease()
	tests.AssertNoError(t, actualErr)
}
