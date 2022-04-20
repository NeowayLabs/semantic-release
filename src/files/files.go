package files

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	e "github.com/NeowayLabs/semantic-release/src/errors"
	style "github.com/NeowayLabs/semantic-release/src/style"
	timeutils "github.com/NeowayLabs/semantic-release/src/time"
)

const (
	changeLogCommitHashLinkFormat = "[%s](https://%s/%s/%s/-/commit/%s)"
)

var (
	versionPyVariable      = "__version__"
	changeLogDefaultFile   = "CHANGELOG.md"
	setupPythonDefaultFile = "setup.py"
)

type FileVersion interface {
	Exists() bool
	UpgradeVersionInSetupPyFile() error
	UpgradeChangelogFile(gitHost, commitType, commitHash, commitMessage, commitAuthor, groupName, projectName string) error
	OpenFile() (*os.File, error)
}

type File struct {
	OutputPath        string
	OriginPath        string
	NewReleaseVersion string
}

type fileVersion struct {
	file File
}

// Exists checks if a file exists within the OS file system.
// Args:
// 		filePath (string): File location in OS file system.
// Returns:
// 		bool: true when file exists, otherwise false.
func (s *fileVersion) Exists() bool {
	if _, err := os.Stat(s.file.OriginPath); err != nil {
		if strings.Contains(err.Error(), e.ErrMsgNoSuchFileOrDirectory) {
			return false
		}
	}

	return true
}

// OpenFile open a file reading it from OS file system.
// Returns:
// 		*os.File: returns a file.
func (s *fileVersion) OpenFile() (*os.File, error) {
	if s.Exists() {
		file, err := os.Open(s.file.OriginPath)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return file, err
	}
	return nil, errors.New(e.ErrMsgNoSuchFileOrDirectory)
}

// UpgradeVersionInSetupPyFile aims to update setup.py file with the new release version.
// It will update the file row containing versionPyVariable. I.e.: __version__.
func (s *fileVersion) UpgradeVersionInSetupPyFile() error {
	defer timeutils.GetElapsedTime("UpgradeVersionInSetupPyFile")()
	log.Printf(style.Yellow+"\nUpgrading version variable in %s file"+style.Reset, s.file.OriginPath)

	var outputData []byte

	file, err := s.OpenFile()
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentLineText := scanner.Text()
		outputLineText := []byte(fmt.Sprintf("%s\n", currentLineText))
		if strings.Contains(currentLineText, fmt.Sprintf("%s =", versionPyVariable)) {
			outputLineText = []byte(fmt.Sprintf("%s = \"%s\"\n", versionPyVariable, s.file.NewReleaseVersion))
		}

		outputData = append(outputData, outputLineText...)
	}

	if err := scanner.Err(); err != nil {
		log.Printf(e.ErrMsgScanningFile, s.file.OriginPath, err)
		return err
	}

	if err = ioutil.WriteFile(s.file.OutputPath, outputData, 0666); err != nil {
		log.Printf(e.ErrMsgWritingUpgradedFile, s.file.OriginPath, s.file.NewReleaseVersion, err)
		return err
	}

	return nil
}

// UpgradeChangelogFile aims to append the new release version with the commit information to the CHANGELOG.md file.
func (s *fileVersion) UpgradeChangelogFile(gitHost, commitType, commitHash, commitMessage, commitAuthor, groupName, projectName string) error {
	log.Printf(style.Yellow+"\nUpgrading %s file"+style.Reset, s.file.OriginPath)
	var outputData []byte

	textToAdd := fmt.Sprintf("\n## v%s\n- %s - %s: %s (%s)\n---\n\n", s.file.NewReleaseVersion,
		commitType,
		fmt.Sprintf(changeLogCommitHashLinkFormat, commitHash[:7], gitHost, groupName, projectName, commitHash),
		commitMessage,
		commitAuthor)
	outputData = append(outputData, []byte(textToAdd)...)

	file, err := s.OpenFile()
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		currentLineText := scanner.Text()
		outputData = append(outputData, []byte(fmt.Sprintf("%s\n", currentLineText))...)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(fmt.Sprintf(e.ErrMsgScanningFile, s.file.OriginPath, err))
		return err
	}

	if err = ioutil.WriteFile(s.file.OutputPath, outputData, 0666); err != nil {
		log.Fatal(fmt.Sprintf(e.ErrMsgWritingUpgradedFile, s.file.OriginPath, s.file.NewReleaseVersion, err))
		return err
	}

	return nil
}

func New(file File) FileVersion {
	return &fileVersion{
		file: file,
	}
}
