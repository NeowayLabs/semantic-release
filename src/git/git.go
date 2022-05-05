package git

import (
	"errors"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const (
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

type Logger interface {
	Info(s string, args ...interface{})
	Error(s string, args ...interface{})
	Warn(s string, args ...interface{})
}

type GitVersioningSystem interface {
	GetChangeHash() (string, error)
	GetChangeAuthorName() (string, error)
	GetChangeAuthorEmail() (string, error)
	GetChangeMessage() (string, error)
	GetCurrentVersion() (string, error)
	UpgradeRemoteRepository(newVersion string) error
}

type ElapsedTime func(functionName string) func()

type GitVersioning struct {
	log                  Logger
	printElapsedTime     ElapsedTime
	url                  string
	destinationDirectory string
	sshKey               string
	replace              bool
	repo                 *git.Repository
}

func (g *GitVersioning) GetChangeHash() (string, error) {
	// TODO: Implement me!
	return "", nil
}
func (g *GitVersioning) GetChangeAuthorName() (string, error) {
	// TODO: Implement me!
	return "", nil
}
func (g *GitVersioning) GetChangeAuthorEmail() (string, error) {
	// TODO: Implement me!
	return "", nil
}
func (g *GitVersioning) GetChangeMessage() (string, error) {
	// TODO: Implement me!
	return "", nil
}
func (g *GitVersioning) GetCurrentVersion() (string, error) {
	// TODO: Implement me!
	return "", nil
}

func (g *GitVersioning) UpgradeRemoteRepository(newVersion string) error {
	// TODO: Implement me!
	return nil
}

// GetPublicKey aims to read ssh key file and returns a go git PublicKey
// Args:
// 		sshKey ([]byte): Byte arry containing the ssh key.
// 		Note: Set up this key in the git settings and pass as a parameter to the CLI tool.
// Returns:
// 		*ssh.PublicKeys: Returns a git public key.
func (g *GitVersioning) getPublicKey(sshKey []byte) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		g.log.Error(err.Error())
		return nil, err
	}
	return publicKey, err
}

// FormatSSHKey aims to properlly format the ssh private key.
// It is necessary once we cannot pass a multiline variable to go build. So it is not possible
// to do docker build --build-arg PRIVATE_KEY=$(PRIVATE_KEY) once this kind of file has several line breaks.
// Args:
// 		lineBreaker (string): Line break holder character. I.e.: #
// Returns:
// 		string: Returns the sshKey string with lineBreaker character replaced by line breaks.
func (g *GitVersioning) formatSSHKey(lineBreaker string) string {
	return strings.Replace(g.sshKey, lineBreaker, "\n", -1)
}

func (g *GitVersioning) validate() error {
	if g.url == "" {
		return errors.New("url cannot be empty")
	}

	if g.destinationDirectory == "" {
		return errors.New("destination directory cannot be empty")
	}

	if g.sshKey == "" {
		return errors.New("ssh key cannot be empty")
	}

	return nil
}

// cloneRepoToDirectory aims to clone the repository from remote to local.
// Args:
// 		auth (*ssh.PublicKeys): Authorization key.
// Returns:
// 		*git.Repository: Returns a repository reference.
// 		err: Error whenever unexpected issues happen.
func (g *GitVersioning) cloneRepoToDirectory() (*git.Repository, error) {
	defer g.printElapsedTime("CloneRepoToDirectory")()

	if err := g.validate(); err != nil {
		g.log.Error(err.Error())
		return nil, err
	}

	sshKey, err := g.getPublicKey([]byte(g.formatSSHKey("#")))
	if err != nil {
		g.log.Error("error while getting public key due to: %s", err)
		return nil, err
	}

	g.log.Info(colorYellow+"cloning repo "+colorCyan+" %s "+colorYellow+" into "+colorCyan+"%s"+colorReset, g.url, g.destinationDirectory)
	repo, err := git.PlainClone(g.destinationDirectory, false, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      g.url,
		Auth:     sshKey,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			g.log.Warn("repository was already cloned")
			if g.replace {
				g.log.Info("removing path %s", g.destinationDirectory)
				err := os.RemoveAll(g.destinationDirectory)
				if err != nil {
					return nil, err
				}
				g.replace = false
				return g.cloneRepoToDirectory()
			}
			return git.PlainOpen(g.destinationDirectory)
		} else {
			g.log.Error("error while cloning gitab repository due to: %s", err)
			return nil, err
		}
	}
	return repo, nil
}

func New(log Logger, printElapsedTime ElapsedTime, url, destinationDirectory, sshKey string, replace bool) (GitVersioningSystem, error) {
	gitLabVersioning := &GitVersioning{
		log:                  log,
		printElapsedTime:     printElapsedTime,
		url:                  url,
		destinationDirectory: destinationDirectory,
		sshKey:               sshKey,
		replace:              replace,
	}

	repo, err := gitLabVersioning.cloneRepoToDirectory()
	if err != nil {
		gitLabVersioning.log.Error("error while cloning repository due to %s", err)
		return nil, err
	}

	gitLabVersioning.repo = repo

	return gitLabVersioning, nil
}
