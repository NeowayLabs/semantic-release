package files

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
	messageTag  = "message:"
)

var (
	changeLogDefaultFile = "CHANGELOG.md"
)

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
	Warn(s string, args ...interface{})
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
	log                Logger
	elapsedTime        ElapsedTime
	versionConrolHost  string
	repositoryRootPath string
	groupName          string
	projectName        string
}

func (f *FileVersion) openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileVersion) findMessageTag(commitMessage string) bool {
	return strings.Contains(strings.ToLower(commitMessage), messageTag)
}

func (f *FileVersion) getMessage(messageRow string) (string, error) {
	startPosition := strings.Index(messageRow, messageTag) + len(messageTag)

	if startPosition-1 == len(messageRow)-1 {
		return "", errors.New("message not found")
	}

	message := strings.TrimSpace(messageRow[startPosition:])
	if strings.ReplaceAll(message, " ", "") == "" {
		return "", errors.New("message not found")
	}

	return message, nil
}

// pretiffy aims to keep a short message based on the commit message, removing extra information such as commit type.
// Args:
// 		commitMessage (string): Full commit message.
// Returns:
// 		string: Returns a commit message with limmited number of characters.
// 		err: Error whenever unexpected issues happen.
func (f *FileVersion) prettifyCommitMessage(commitMessage string) (string, error) {

	var message string
	splitedMessage := strings.Split(commitMessage, "\n")

	for _, row := range splitedMessage {
		row := strings.ToLower(row)
		if f.findMessageTag(row) {

			currentMessage, err := f.getMessage(row)
			if err != nil {
				return "", err
			}
			message = currentMessage
		}
	}

	if message == "" {
		return "", errors.New("commit message has no tag 'message:'")
	}
	// Limmit message to 150 characters to avoid long messages.
	if len(message) >= 150 {
		message = fmt.Sprintf("%s...", message[:150])
	}

	// Upper only first letter of the sentence.
	message = fmt.Sprintf("%s%s", strings.ToUpper(message[:1]), message[1:])
	return message, nil
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

// UpgradeVariableInFiles aims to update given files with the new release version.
// It will update the files row containing a given variable name.
// I.e.:
// err := UpgradeVariableInFiles(UpgradeFiles{Files: []UpgradeFile{{Path: "./setup.py", DestinationPath: "", VariableName: "__version__"}}), "1.0.1")
//  	From: __version__ = 1.0.0
//  	To:   __version__ = 1.0.1
func (f *FileVersion) UpgradeVariableInFiles(filesToUpgrade interface{}, newVersion string) error {
	defer f.elapsedTime("UpgradeVariableInFiles")()

	var outputData []byte

	filesToUpdate, err := f.unmarshalUpgradeFiles(filesToUpgrade)
	if err != nil {
		return err
	}

	for _, currentFile := range filesToUpdate.Files {
		f.log.Info(colorYellow+"Upgrading version variable in %s file"+colorReset, currentFile.Path)

		file, err := f.openFile(currentFile.Path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		variableNameFound := false
		for scanner.Scan() {
			currentLineText := scanner.Text()
			outputLineText := []byte(fmt.Sprintf("%s\n", currentLineText))

			if f.containsVariableNameInText(currentLineText, currentFile.VariableName, ":=") {
				outputLineText = []byte(fmt.Sprintf("%s := \"%s\"\n", currentFile.VariableName, newVersion))
				variableNameFound = true
			}

			if f.containsVariableNameInText(currentLineText, currentFile.VariableName, "=") {
				outputLineText = []byte(fmt.Sprintf("%s = \"%s\"\n", currentFile.VariableName, newVersion))
				variableNameFound = true
			}

			outputData = append(outputData, outputLineText...)
		}

		if !variableNameFound {
			return fmt.Errorf("variable name `%s` not found on file `%s`", currentFile.VariableName, currentFile.Path)
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		destinationPath := f.setDefaultPath(currentFile.DestinationPath, currentFile.Path)
		if err = ioutil.WriteFile(destinationPath, outputData, 0666); err != nil {
			return err
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

// UpgradeChangelog aims to append the new release version with the commit information to the CHANGELOG.md file.
func (f *FileVersion) UpgradeChangeLog(path, destinationPath string, chageLogInfo interface{}) error {
	defer f.elapsedTime("UpgradeChangeLog")()

	originPath := f.setDefaultPath(path, fmt.Sprintf("%s/%s", f.repositoryRootPath, changeLogDefaultFile))

	f.log.Info(colorYellow+"Upgrading %s file"+colorReset, originPath)

	changelog, err := f.unmarshalChangesInfo(chageLogInfo)
	if err != nil {
		return err
	}

	if err := f.validateChangesInfo(*changelog); err != nil {
		return err
	}

	commitMessage, err := f.prettifyCommitMessage(changelog.Message)
	if err != nil {
		return err
	}

	// textToAdd
	// I.e.:
	// ## v1.0.0:
	// - feat - [b25a9af](https://gilabhost/groupName/projectName/commit/b25a9af78c30de0d03ca2ee6d18c66bbc4804395): Commit message here (@user.name)
	textToAdd := fmt.Sprintf("\n## v%s\n- %s - %s: %s (%s)\n---\n\n",
		changelog.NewVersion,
		changelog.ChangeType,
		f.getCommitUrl(*changelog),
		commitMessage,
		f.prettifyEmail(changelog.AuthorEmail))

	var outputData []byte
	outputData = append(outputData, []byte(textToAdd)...)

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

	destinationPath = f.setDefaultPath(destinationPath, path)
	if err = ioutil.WriteFile(destinationPath, outputData, 0666); err != nil {
		return fmt.Errorf("\n\nerror while writing %s file with new version %s due to: %w", originPath, changelog.NewVersion, err)
	}

	return nil
}

func New(log Logger, elapsedTime ElapsedTime, versionConrolHost, repositoryRootPath, groupName, projectName string) *FileVersion {
	return &FileVersion{
		log:                log,
		elapsedTime:        elapsedTime,
		versionConrolHost:  versionConrolHost,
		repositoryRootPath: repositoryRootPath,
		groupName:          groupName,
		projectName:        projectName,
	}
}
