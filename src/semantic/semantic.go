package semantic

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorBGRed  = "\033[41;1;37m"
)

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
}

type RepositoryVersionControl interface {
	GetChangeHash() string
	GetChangeAuthorName() string
	GetChangeAuthorEmail() string
	GetChangeMessage() string
	GetCurrentVersion() string
	UpgradeRemoteRepository(newVersion string) error
	GetCommitHistory() []*object.Commit
	GetCommitHistoryDiff() []*object.Commit
}

type VersionControl interface {
	GetCommitChangeType(commitMessage string) (string, error)
	GetNewVersion(commitMessage string, currentVersion string) (string, error)
	MustSkipVersioning(commitMessage string) bool
}

type FilesVersionControl interface {
	UpgradeChangeLog(path, destinationPath string, chageLogInfo interface{}) error
	UpgradeVariableInFiles(filesToUpgrade interface{}, newVersion string) error
}

type ChangesInfo struct {
	Hash           string
	AuthorName     string
	AuthorEmail    string
	Message        string
	CurrentVersion string
	NewVersion     string
	ChangeType     string
}

type Semantic struct {
	log                   Logger
	rootPath              string
	filesToUpdateVariable interface{}
	repoVersionControl    RepositoryVersionControl
	versionControl        VersionControl
	filesVersionControl   FilesVersionControl
}

func (s *Semantic) GenerateNewRelease() error {
	changesInfo := &ChangesInfo{
		Hash:           s.repoVersionControl.GetChangeHash(),
		AuthorName:     s.repoVersionControl.GetChangeAuthorName(),
		AuthorEmail:    s.repoVersionControl.GetChangeAuthorEmail(),
		Message:        s.repoVersionControl.GetChangeMessage(),
		CurrentVersion: s.repoVersionControl.GetCurrentVersion(),
	}

	if s.versionControl.MustSkipVersioning(changesInfo.Message) {
		s.log.Info(colorCyan + "Semantic Release has been skiped by commit message tag [skip]" + colorReset)
		return nil
	}

	newVersion, err := s.versionControl.GetNewVersion(changesInfo.Message, changesInfo.CurrentVersion)
	if err != nil {
		return errors.New("error while getting new version due to: " + err.Error())
	}

	changesInfo.NewVersion = newVersion

	commitChangeType, err := s.versionControl.GetCommitChangeType(changesInfo.Message)
	if err != nil {
		return fmt.Errorf("error while getting commit change type due to: %s", err.Error())
	}

	changesInfo.ChangeType = commitChangeType

	s.log.Info(colorBGRed + "MOST RECENT COMMIT:" + colorReset)
	s.log.Info("Hash: %s", changesInfo.Hash)
	s.log.Info("Author Name: %s", changesInfo.AuthorName)
	s.log.Info("Author Email: %s", changesInfo.AuthorEmail)
	s.log.Info("Message: %s", changesInfo.Message)
	s.log.Info("Current Version: %s", changesInfo.CurrentVersion)
	s.log.Info(fmt.Sprintf("Commit change type: "+colorYellow+"%s"+colorReset, commitChangeType))
	s.log.Info("New Version: %s", changesInfo.NewVersion)

	if err := s.filesVersionControl.UpgradeChangeLog(fmt.Sprintf("%s/CHANGELOG.md", s.rootPath), "", changesInfo); err != nil {
		return errors.New("error while upgrading changelog file due to: " + err.Error())
	}

	if s.filesToUpdateVariable != nil {
		if err := s.filesVersionControl.UpgradeVariableInFiles(s.filesToUpdateVariable, changesInfo.NewVersion); err != nil {
			return errors.New("error while upgrading variables in files due to: " + err.Error())
		}
	}

	if err := s.repoVersionControl.UpgradeRemoteRepository(newVersion); err != nil {
		return errors.New("error while upgrading remote repository due to: " + err.Error())
	}

	return nil
}

func (s *Semantic) isValidMessage(message string) bool {
	_, err := s.versionControl.GetCommitChangeType(message)
	if err != nil {
		if err.Error() == "change type not found" {
			s.log.Error("change type not found")
		}
		return false
	}

	return strings.Contains(strings.ToLower(message), "message:")
}

func (s *Semantic) CommitLint() error {
	commitHistoryDiff := s.repoVersionControl.GetCommitHistoryDiff()

	areThereWrongCommits := false
	for _, commit := range commitHistoryDiff {
		if !s.isValidMessage(commit.Message) {
			s.log.Error(colorRed+"commit message "+colorYellow+"( %s )"+colorRed+" does not meet semantic-release pattern "+colorYellow+"( type: [commit type], message: message here.)"+colorReset, strings.TrimSuffix(commit.Message, "\n"))
			areThereWrongCommits = true
		}
	}
	if areThereWrongCommits {
		s.log.Error(colorRed + "You can use " + colorBGRed + "git rebase -i HEAD~n" + colorReset + colorRed + " and edit the commit list with reword before each commit message." + colorReset)
		return errors.New("commit messages dos not meet semantic-release pattern")
	}

	return nil
}

func New(log Logger, rootPath string, filesToUpdateVariable interface{}, repoVersionControl RepositoryVersionControl, filesVersionControl FilesVersionControl, versionControl VersionControl) *Semantic {
	return &Semantic{
		log:                   log,
		rootPath:              rootPath,
		filesToUpdateVariable: filesToUpdateVariable,
		repoVersionControl:    repoVersionControl,
		filesVersionControl:   filesVersionControl,
		versionControl:        versionControl,
	}
}
