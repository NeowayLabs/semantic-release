//go:build unit
// +build unit

package files_test

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/files"
	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

type ChangesInfoMock struct {
	Hash           string
	AuthorName     string
	AuthorEmail    string
	Message        string
	CurrentVersion string
	NewVersion     string
	ChangeType     string
}

type UpgradeFilesMock struct {
	Files []UpgradeFileMock
}

type UpgradeFileMock struct {
	Path            string
	DestinationPath string
	VariableName    string
}

func printElapsedTimeMock(functionName string) func() {
	return func() {
		fmt.Printf("%s done.", functionName)
	}
}

type fixture struct {
	log                *log.Log
	versionControlHost string
	repositoryRootPath string
	groupName          string
	projectName        string
}

func setup(t *testing.T) *fixture {
	logger, err := log.New("test", "1.0.0", "debug")
	if err != nil {
		t.Errorf("error while getting log due to %s", err.Error())
	}

	return &fixture{log: logger}
}

func (f *fixture) newFiles() *files.FileVersion {
	return files.New(f.log, printElapsedTimeMock, f.versionControlHost, f.repositoryRootPath, f.groupName, f.projectName)
}

func TestUpgradeVariableInFilesNoError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	filesToUpgrade := UpgradeFilesMock{Files: []UpgradeFileMock{{Path: "mock/setup_mock.py", VariableName: "__version__"}}}
	err := filesVersion.UpgradeVariableInFiles(filesToUpgrade, "1.0.1")
	tests.AssertNoError(t, err)

}

func TestUpgradeVariableInFilesVariableNameNotFoundError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	filesToUpgrade := UpgradeFilesMock{Files: []UpgradeFileMock{{Path: "mock/setup_mock.py", VariableName: "version"}}}
	err := filesVersion.UpgradeVariableInFiles(filesToUpgrade, "1.0.1")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while getting file output data due to: variable name `version` not found on file `mock/setup_mock.py`", err.Error())

}

func TestUpgradeVariableInFilesMarsahlError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	err := filesVersion.UpgradeVariableInFiles(make(chan int), "1.0.1")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error unmarshalling files to upgrade due to: error marshalling files to uptade", err.Error())

}

func TestUpgradeVariableInFilesUnmarsahlError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	err := filesVersion.UpgradeVariableInFiles("", "1.0.1")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error unmarshalling files to upgrade due to: error unmarshalling files to uptade", err.Error())

}

func TestUpgradeVariableInFilesOpenFileError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	filesToUpgrade := UpgradeFilesMock{Files: []UpgradeFileMock{{Path: "mock/setup_404.py", VariableName: "__version__"}}}
	err := filesVersion.UpgradeVariableInFiles(filesToUpgrade, "1.0.1")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while oppening file due to: open mock/setup_404.py: no such file or directory", err.Error())

}

func TestUpgradeVariableInFilesWriteFileError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	filesToUpgrade := UpgradeFilesMock{Files: []UpgradeFileMock{{Path: "mock/setup_mock.py", DestinationPath: "mock/test/setup_mock_404.py", VariableName: "__version__"}}}
	err := filesVersion.UpgradeVariableInFiles(filesToUpgrade, "1.0.1")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while writing upgrade variables in file due to: error while writing file mock/test/setup_mock_404.py due to: open mock/test/setup_mock_404.py: no such file or directory", err.Error())

}

func TestUpgradeChangeLogNoError(t *testing.T) {
	f := setup(t)
	f.versionControlHost = "gitlab.com"
	f.groupName = "dataplatform"
	f.projectName = "test"
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		Hash:           "b25a9af78c30de0d03ca2ee6d18c66bbc4804395",
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		Message:        "type: [feat], Message: This is a short message to write to CHANGELOG.md file.",
		CurrentVersion: "1.0.1",
		NewVersion:     "1.1.0",
		ChangeType:     "feat",
	}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertNoError(t, err)
}

func TestUpgradeChangeLogLongMessageCutNoError(t *testing.T) {
	f := setup(t)
	f.versionControlHost = "gitlab.com"
	f.groupName = "dataplatform"
	f.projectName = "test"
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		Hash:           "b25a9af78c30de0d03ca2ee6d18c66bbc4804395",
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		Message:        "type: [feat], Message: This is a long message to write to CHANGELOG.md file. Bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo bar foo cut here.",
		CurrentVersion: "1.0.1",
		NewVersion:     "1.1.0",
		ChangeType:     "feat",
	}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertNoError(t, err)
}

func TestUpgradeChangeLogMarshalChangeLogInfoError(t *testing.T) {
	f := setup(t)
	f.versionControlHost = "gitlab.com"
	f.groupName = "dataplatform"
	f.projectName = "test"
	filesVersion := f.newFiles()

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", make(chan int))
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error unmarshalling changes info due to: error marshalling files to changelog information", err.Error())
}

func TestUpgradeChangeLogUnmarshalError(t *testing.T) {
	f := setup(t)
	f.versionControlHost = "gitlab.com"
	f.groupName = "dataplatform"
	f.projectName = "test"
	filesVersion := f.newFiles()

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", "")
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error unmarshalling changes info due to: error unmarshalling changelog information", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyAuthorNameError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: author name cannot be empty", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyAuthorEmailBadEntryError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{AuthorName: "Administrator", AuthorEmail: "admingitlab.com"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: bad author email entry", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyAuthorEmailError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{AuthorName: "Administrator"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: bad author email entry", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyChangeTypeError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{AuthorName: "Administrator", AuthorEmail: "admin@git.com"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: change type cannot be empty", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoHashError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{AuthorName: "Administrator", AuthorEmail: "admin@git.com", ChangeType: "feat", Hash: "b25a9a"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: hash string must have at least 7 characters", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyMessageError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{AuthorName: "Administrator", AuthorEmail: "admin@git.com", ChangeType: "feat", Hash: "b25a9af"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: message cannot be empty", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyCurrentVersionError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{AuthorName: "Administrator", AuthorEmail: "admin@git.com", ChangeType: "feat", Hash: "b25a9af", Message: "anything"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: current version cannot be empty", err.Error())
}

func TestUpgradeChangeLogValidateChangesInfoEmptyNewVersionError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		ChangeType:     "feat",
		Hash:           "b25a9af",
		Message:        "anything",
		CurrentVersion: "1.0.0"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error validating changelog info due to: new version cannot be empty", err.Error())
}

func TestUpgradeChangeLogMessageNotFoundError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		ChangeType:     "feat",
		Hash:           "b25a9af",
		Message:        "message:",
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while formatting changelog content due to: prettify commit message error: error while getting message due to: message not found", err.Error())
}

func TestUpgradeChangeLogTagMessageNotFoundEmptyError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		ChangeType:     "feat",
		Hash:           "b25a9af",
		Message:        "message:  ",
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while formatting changelog content due to: prettify commit message error: error while getting message due to: message not found", err.Error())
}

func TestUpgradeChangeLogTagMessageNotFoundError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		ChangeType:     "feat",
		Hash:           "b25a9af",
		Message:        "type: [feat]",
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while formatting changelog content due to: prettify commit message error: commit message has no tag 'message:'", err.Error())
}

func TestUpgradeChangeLogOpenFileError(t *testing.T) {
	f := setup(t)
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		ChangeType:     "feat",
		Hash:           "b25a9af",
		Message:        "type: [feat], message: Test.",
		CurrentVersion: "1.0.0",
		NewVersion:     "1.1.0"}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG.md", "", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while openning changelog file due to: error while oppening file due to: open mock/CHANGELOG.md: no such file or directory", err.Error())
}

func TestUpgradeChangeLogWriteFileError(t *testing.T) {
	f := setup(t)
	f.versionControlHost = "gitlab.com"
	f.groupName = "dataplatform"
	f.projectName = "test"
	filesVersion := f.newFiles()

	changelog := ChangesInfoMock{
		Hash:           "b25a9af78c30de0d03ca2ee6d18c66bbc4804395",
		AuthorName:     "Administrator",
		AuthorEmail:    "admin@git.com",
		Message:        "type: [feat], Message: This is a short message to write to CHANGELOG.md file.",
		CurrentVersion: "1.0.1",
		NewVersion:     "1.1.0",
		ChangeType:     "feat",
	}

	err := filesVersion.UpgradeChangeLog("mock/CHANGELOG_MOCK.md", "mock/test/CHANGELOG_404.md", changelog)
	tests.AssertError(t, err)
	tests.AssertEqualValues(t, "error while writing new version to changelog file due to: error while writing file mock/test/CHANGELOG_404.md due to: open mock/test/CHANGELOG_404.md: no such file or directory", err.Error())
}
