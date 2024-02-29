package files

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

var (
	changeLogDefaultFile = "CHANGELOG.md"
)

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
	Warn(s string, args ...interface{})
}

type CommitMessageManager interface {
	PrettifyCommitMessage(commitMessage string) (string, error)
}

type ElapsedTime func(functionName string) func()

type ChangesInfo struct {
	Hash           string
	AuthorName     string
	AuthorEmail    string
	Message        string
	CurrentVersion string
	NewVersion     string
	ChangeType     string
}

type UpgradeFiles struct {
	Files []UpgradeFile
}

type UpgradeFile struct {
	Path            string
	DestinationPath string
	VariableName    string
}

type FileVersion struct {
	log                  Logger
	elapsedTime          ElapsedTime
	versionConrolHost    string
	repositoryRootPath   string
	groupName            string
	projectName          string
	variableNameFound    bool
	commitMessageManager CommitMessageManager
}

func (f *FileVersion) openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while oppening file due to: %w", err)
	}
	return file, nil
}

func (f *FileVersion) containsVariableNameInText(text, variableName, copySignal string) bool {
	text = strings.ReplaceAll(text, " ", "")
	return strings.Contains(text, fmt.Sprintf("%s%s", variableName, copySignal))
}

func (f *FileVersion) setDefaultPath(path, newPath string) string {
	if path != "" {
		return path
	}
	return newPath
}

func (f *FileVersion) unmarshalUpgradeFiles(filesToUpgrade interface{}) (*UpgradeFiles, error) {
	filesToUpdateBytes, err := json.Marshal(filesToUpgrade)
	if err != nil {
		return nil, errors.New("error marshalling files to uptade")
	}

	var filesToUpdateResult UpgradeFiles
	if err := json.Unmarshal(filesToUpdateBytes, &filesToUpdateResult); err != nil {
		return nil, errors.New("error unmarshalling files to uptade")
	}

	return &filesToUpdateResult, nil
}

func (f *FileVersion) writeFile(destinationPath, originPath string, content []byte) error {
	destination := f.setDefaultPath(destinationPath, originPath)

	if err := os.WriteFile(destination, content, 0666); err != nil {
		return fmt.Errorf("error while writing file %s due to: %w", destination, err)
	}

	return nil
}

func (f *FileVersion) getFileOutputContent(scanner *bufio.Scanner, file UpgradeFile, newVersion string) ([]byte, error) {
	var outputData []byte

	f.variableNameFound = false
	for scanner.Scan() {
		currentLineText := scanner.Text()
		outputLineText := []byte(fmt.Sprintf("%s\n", currentLineText))

		if f.containsVariableNameInText(currentLineText, file.VariableName, ":=") {
			outputLineText = []byte(fmt.Sprintf("%s := \"%s\"\n", file.VariableName, newVersion))
			f.variableNameFound = true
		}

		if f.containsVariableNameInText(currentLineText, file.VariableName, "=") {
			outputLineText = []byte(fmt.Sprintf("%s = \"%s\"\n", file.VariableName, newVersion))
			f.variableNameFound = true
		}

		outputData = append(outputData, outputLineText...)
	}

	if !f.variableNameFound {
		return nil, fmt.Errorf("variable name `%s` not found on file `%s`", file.VariableName, file.Path)
	}

	return outputData, nil
}

// UpgradeVariableInFiles aims to update given files with the new release version.
// It will update the files row containing a given variable name.
// I.e.:
// err := UpgradeVariableInFiles(UpgradeFiles{Files: []UpgradeFile{{Path: "./setup.py", DestinationPath: "", VariableName: "__version__"}}), "1.0.1")
//
//	From: __version__ = 1.0.0
//	To:   __version__ = 1.0.1
func (f *FileVersion) UpgradeVariableInFiles(filesToUpgrade interface{}, newVersion string) error {
	defer f.elapsedTime("UpgradeVariableInFiles")()

	filesToUpdate, err := f.unmarshalUpgradeFiles(filesToUpgrade)
	if err != nil {
		return fmt.Errorf("error unmarshalling files to upgrade due to: %w", err)
	}

	for _, currentFile := range filesToUpdate.Files {
		f.log.Info(colorYellow+"Upgrading version variable in %s file"+colorReset, currentFile.Path)

		file, err := f.openFile(currentFile.Path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		outputData, err := f.getFileOutputContent(scanner, currentFile, newVersion)
		if err != nil {
			return fmt.Errorf("error while getting file output data due to: %w", err)
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error while scanning file %s due to: %w", currentFile.Path, err)
		}

		if err = f.writeFile(currentFile.DestinationPath, currentFile.Path, outputData); err != nil {
			return fmt.Errorf("error while writing upgrade variables in file due to: %w", err)
		}
	}

	return nil
}

func (f *FileVersion) validateChangesInfo(changelog ChangesInfo) error {

	if changelog.AuthorName == "" {
		return errors.New("author name cannot be empty")
	}

	if changelog.AuthorEmail == "" || !strings.Contains(changelog.AuthorEmail, "@") {
		return errors.New("bad author email entry")
	}

	if changelog.ChangeType == "" {
		return errors.New("change type cannot be empty")
	}

	if len(changelog.Hash) < 7 {
		return errors.New("hash string must have at least 7 characters")
	}

	if changelog.Message == "" {
		return errors.New("message cannot be empty")
	}

	if changelog.CurrentVersion == "" {
		return errors.New("current version cannot be empty")
	}

	if changelog.NewVersion == "" {
		return errors.New("new version cannot be empty")
	}

	return nil
}

func (f *FileVersion) abbreviateHash(hash string) string {
	return hash[:7]
}

func (f *FileVersion) getCommitUrl(changelog ChangesInfo) string {
	return fmt.Sprintf("[%s](https://%s/%s/%s/commit/%s)", f.abbreviateHash(changelog.Hash), f.versionConrolHost, f.groupName, f.projectName, changelog.Hash)
}

func (f *FileVersion) prettifyEmail(email string) string {
	splitedEmail := strings.Split(email, "@")
	return fmt.Sprintf("@%s", splitedEmail[0])
}

func (f *FileVersion) unmarshalChangesInfo(changes interface{}) (*ChangesInfo, error) {
	changeLogInfoBytes, err := json.Marshal(changes)
	if err != nil {
		return nil, errors.New("error marshalling files to changelog information")
	}

	var changelog ChangesInfo
	if err := json.Unmarshal(changeLogInfoBytes, &changelog); err != nil {
		return nil, errors.New("error unmarshalling changelog information")
	}

	return &changelog, nil
}

func (f *FileVersion) formatChangeLogContent(changes *ChangesInfo) (string, error) {
	commitMessage, err := f.commitMessageManager.PrettifyCommitMessage(changes.Message)
	if err != nil {
		return "", fmt.Errorf("prettify commit message error: %w", err)
	}

	// textToAdd
	// I.e.:
	// ## v1.0.0:
	// - feat - [b25a9af](https://gilabhost/groupName/projectName/commit/b25a9af78c30de0d03ca2ee6d18c66bbc4804395): Commit message here (@user.name)
	return fmt.Sprintf("\n## v%s\n- %s - %s: %s (%s)\n---\n\n",
		changes.NewVersion,
		changes.ChangeType,
		f.getCommitUrl(*changes),
		commitMessage,
		f.prettifyEmail(changes.AuthorEmail)), nil
}

// UpgradeChangelog aims to append the new release version with the commit information to the CHANGELOG.md file.
func (f *FileVersion) UpgradeChangeLog(path, destinationPath string, chageLogInfo interface{}) error {
	defer f.elapsedTime("UpgradeChangeLog")()

	originPath := f.setDefaultPath(path, fmt.Sprintf("%s/%s", f.repositoryRootPath, changeLogDefaultFile))

	f.log.Info(colorYellow+"Upgrading %s file"+colorReset, originPath)

	changelog, err := f.unmarshalChangesInfo(chageLogInfo)
	if err != nil {
		return fmt.Errorf("error unmarshalling changes info due to: %w", err)
	}

	if err := f.validateChangesInfo(*changelog); err != nil {
		return fmt.Errorf("error validating changelog info due to: %w", err)
	}

	textToAdd, err := f.formatChangeLogContent(changelog)
	if err != nil {
		return fmt.Errorf("error while formatting changelog content due to: %w", err)
	}

	outputData := []byte(textToAdd)

	file, err := f.openFile(originPath)
	if err != nil {
		return fmt.Errorf("error while openning changelog file due to: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentLineText := scanner.Text()
		outputData = append(outputData, []byte(fmt.Sprintf("%s\n", currentLineText))...)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("\n\nerror while scanning file: %s due to: %w", originPath, err)
	}

	if err = f.writeFile(destinationPath, path, outputData); err != nil {
		return fmt.Errorf("error while writing new version to changelog file due to: %w", err)
	}

	return nil
}

func New(log Logger, elapsedTime ElapsedTime, versionConrolHost, repositoryRootPath, groupName, projectName string, commitMessageManager CommitMessageManager) *FileVersion {
	return &FileVersion{
		log:                  log,
		elapsedTime:          elapsedTime,
		versionConrolHost:    versionConrolHost,
		repositoryRootPath:   repositoryRootPath,
		groupName:            groupName,
		projectName:          projectName,
		commitMessageManager: commitMessageManager,
	}
}
