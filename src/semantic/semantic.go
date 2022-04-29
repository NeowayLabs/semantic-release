package semantic

import (
	"encoding/json"
	"errors"
)

type RepositoryVersionControl interface {
	GetChangeHash() (string, error)
	GetChangeAuthorName() (string, error)
	GetChangeAuthorEmail() (string, error)
	GetChangeMessage() (string, error)
	GetCurrentVersion() (string, error)
	UpgradeRemoteRepository(newVersion string) error
}

type VersionControl interface {
	GetNewVersion(commitMessage string, currentVersion string) (string, error)
}

type FilesVersionControl interface {
	UpgradeChangeLog(path string, chagelogInfo interface{}, newVersion string) error
	UpgradeVariableInFiles(filesInfo interface{}, newVersion string) error
}

type upgradeFiles struct {
	files []upgradeFile
}

type upgradeFile struct {
	path         string
	variableName string
}

type changesInfo struct {
	hash           string
	authorName     string
	authorEmail    string
	message        string
	currentVersion string
}

type Semantic struct {
	rootPath              string
	filesToUpdateVariable interface{}
	repoVersionControl    RepositoryVersionControl
	versionControl        VersionControl
	filesVersionControl   FilesVersionControl
}

func (s *Semantic) getChangesInformation() (*changesInfo, error) {
	hash, err := s.repoVersionControl.GetChangeHash()
	if err != nil {
		return nil, errors.New("error getting hash: " + err.Error())
	}

	authorName, err := s.repoVersionControl.GetChangeAuthorName()
	if err != nil {
		return nil, errors.New("error getting author name: " + err.Error())
	}

	authorEmail, err := s.repoVersionControl.GetChangeAuthorEmail()
	if err != nil {
		return nil, errors.New("error getting author email: " + err.Error())
	}

	message, err := s.repoVersionControl.GetChangeMessage()
	if err != nil {
		return nil, errors.New("error getting message: " + err.Error())
	}

	currentVersion, err := s.repoVersionControl.GetCurrentVersion()
	if err != nil {
		return nil, errors.New("error getting current version: " + err.Error())
	}

	return &changesInfo{hash: hash,
		authorName:     authorName,
		authorEmail:    authorEmail,
		message:        message,
		currentVersion: currentVersion}, nil
}

func (s *Semantic) GenerateNewRelease() error {
	changesInfo, err := s.getChangesInformation()
	if err != nil {
		return errors.New("error while getting changes information due to: " + err.Error())
	}

	newVersion, err := s.versionControl.GetNewVersion(changesInfo.message, changesInfo.currentVersion)
	if err != nil {
		return errors.New("error while getting new version due to: " + err.Error())
	}

	if err := s.filesVersionControl.UpgradeChangeLog(s.rootPath, changesInfo, newVersion); err != nil {
		return errors.New("error while upgrading changelog file due to: " + err.Error())
	}

	if s.filesToUpdateVariable != nil {

		filesToUpdateBytes, err := json.Marshal(s.filesToUpdateVariable)
		if err != nil {
			return errors.New("error marshalling files to uptade information")
		}

		var filesToUpdateVariable upgradeFiles
		if err := json.Unmarshal(filesToUpdateBytes, &filesToUpdateVariable); err != nil {
			return errors.New("error unmarshalling files to uptade information")
		}

		if err := s.filesVersionControl.UpgradeVariableInFiles(s.filesToUpdateVariable, newVersion); err != nil {
			return errors.New("error while upgrading variables in files due to: " + err.Error())
		}
	}

	if err := s.repoVersionControl.UpgradeRemoteRepository(newVersion); err != nil {
		return errors.New("error while upgrading remote repository due to: " + err.Error())
	}

	return nil
}

func New(rootPath string, filesToUpdateVariable interface{}, repoVersionControl RepositoryVersionControl, filesVersionControl FilesVersionControl, versionControl VersionControl) *Semantic {
	return &Semantic{
		rootPath:              rootPath,
		filesToUpdateVariable: filesToUpdateVariable,
		repoVersionControl:    repoVersionControl,
		filesVersionControl:   filesVersionControl,
		versionControl:        versionControl,
	}
}
