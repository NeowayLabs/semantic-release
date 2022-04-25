package version

import (
	"errors"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/tests"
)

const (
	shouldbeEqualMessage = "Values should be equal:\n Actual: %v\nExpected: %v"
)

type RepositoryVersionControlMock struct {
	message string
	err     error
}

func (r *RepositoryVersionControlMock) GetMessage() (string, error) {
	return r.message, r.err
}

func (r *RepositoryVersionControlMock) GetVersionUpdateType(message string) (string, error) {
	return message, r.err
}

func (r *RepositoryVersionControlMock) UpgradeRemoteRepository() error {
	return r.err
}

type FilesVersionControlMock struct {
	err error
}

func (f *FilesVersionControlMock) UpgradeChangeLog(path string) error {
	return f.err
}
func (f *FilesVersionControlMock) UpgradeVariableInFiles(variableName string, filesList []string) error {
	return f.err
}

type fixture struct {
	repoVersionMock  *RepositoryVersionControlMock
	filesVersionMock *FilesVersionControlMock
}

func setup() *fixture {
	return &fixture{repoVersionMock: &RepositoryVersionControlMock{}, filesVersionMock: &FilesVersionControlMock{}}
}

func TestGetMessageError(t *testing.T) {
	f := setup()
	f.repoVersionMock.err = errors.New("Something went wrong while getting message.")
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	_, actualErr := versionService.repoVersionControl.GetMessage()

	tests.AssertError(t, actualErr)
}

func TestGetMessageNoError(t *testing.T) {
	f := setup()
	f.repoVersionMock.message = "This is a message"
	expectedMessage := "This is a message"

	versionService := NewService(f.repoVersionMock, f.filesVersionMock)

	actualMessage, actualErr := versionService.repoVersionControl.GetMessage()
	tests.AssertNoError(t, actualErr)

	if expectedMessage != actualMessage {
		t.Errorf(shouldbeEqualMessage, actualMessage, expectedMessage)
	}
}

func TestGetVersionUpdateTypeError(t *testing.T) {
	f := setup()
	f.repoVersionMock.err = errors.New("Something went wrong while getting update type.")
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)

	_, actualErr := versionService.repoVersionControl.GetVersionUpdateType("")

	tests.AssertError(t, actualErr)
}

func TestGetVersionUpdateNoError(t *testing.T) {
	f := setup()
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	expectedType := "feat"
	actualType, actualErr := versionService.repoVersionControl.GetVersionUpdateType("feat")

	tests.AssertNoError(t, actualErr)
	if expectedType != actualType {
		t.Errorf(shouldbeEqualMessage, actualType, expectedType)
	}
}

func TestUpgradeRemoteRepositoryError(t *testing.T) {
	f := setup()
	f.repoVersionMock.err = errors.New("Something went wrong while upgrading remote repository.")
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	actualErr := versionService.repoVersionControl.UpgradeRemoteRepository()

	tests.AssertError(t, actualErr)
}

func TestUpgradeRemoteRepositoryNoError(t *testing.T) {
	f := setup()
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	actualErr := versionService.repoVersionControl.UpgradeRemoteRepository()

	tests.AssertNoError(t, actualErr)
}

func TestUpgradeChangeLogError(t *testing.T) {
	f := setup()
	f.filesVersionMock.err = errors.New("Something went wrong while upgrading changelog.")
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	actualErr := versionService.filesVersionControl.UpgradeChangeLog("")

	tests.AssertError(t, actualErr)
}

func TestUpgradeChangeLogNoError(t *testing.T) {
	f := setup()
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	actualErr := versionService.filesVersionControl.UpgradeChangeLog("")

	tests.AssertNoError(t, actualErr)
}

func TestUpgradeVariableInFilesError(t *testing.T) {
	f := setup()
	f.filesVersionMock.err = errors.New("Something went wrong while upgrading files variable.")
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	actualErr := versionService.filesVersionControl.UpgradeVariableInFiles("", []string{})

	tests.AssertError(t, actualErr)
}

func TestUpgradeVariableInFilesNoError(t *testing.T) {
	f := setup()
	versionService := NewService(f.repoVersionMock, f.filesVersionMock)
	actualErr := versionService.filesVersionControl.UpgradeVariableInFiles("", []string{})

	tests.AssertNoError(t, actualErr)
}
