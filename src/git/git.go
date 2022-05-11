package git

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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

type ElapsedTime func(functionName string) func()

type GitVersioning struct {
	log                  Logger
	printElapsedTime     ElapsedTime
	url                  string
	destinationDirectory string
	repo                 *git.Repository
	username             string
	password             string
}

func (g *GitVersioning) validate() error {
	if g.url == "" {
		return errors.New("url cannot be empty")
	}

	if g.destinationDirectory == "" {
		return errors.New("destination directory cannot be empty")
	}

	if g.username == "" {
		return errors.New("username cannot be empty")
	}

	if g.password == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}

// cloneRepoToDirectory aims to clone the repository from remote to local.
// Returns:
// 		*git.Repository: Returns a repository reference.
// 		err: Error whenever unexpected issues happen.
func (g *GitVersioning) cloneRepoToDirectory() (*git.Repository, error) {
	defer g.printElapsedTime("CloneRepoToDirectory")()

	g.log.Info(colorYellow+"cloning repo "+colorCyan+" %s "+colorYellow+" into "+colorCyan+"%s"+colorReset, g.url, g.destinationDirectory)
	repo, err := git.PlainClone(g.destinationDirectory, false, &git.CloneOptions{
		Progress: os.Stdout,
		URL:      g.url,
		Auth: &http.BasicAuth{Username: g.username,
			Password: g.password,
		},
		InsecureSkipTLS: true,
	})

	if err == nil {
		return repo, nil
	}

	if err == git.ErrRepositoryAlreadyExists {
		g.log.Warn("repository was already cloned")
		return git.PlainOpen(g.destinationDirectory)
	}
	g.log.Error("error while cloning gitab repository due to: %s", err)
	return nil, err
}

func New(log Logger, printElapsedTime ElapsedTime, url, username, password, destinationDirectory string) (*GitVersioning, error) {
	gitLabVersioning := &GitVersioning{
		log:                  log,
		printElapsedTime:     printElapsedTime,
		username:             username,
		password:             password,
		url:                  url,
		destinationDirectory: destinationDirectory,
	}

	if err := gitLabVersioning.validate(); err != nil {
		gitLabVersioning.log.Error(err.Error())
		return nil, err
	}

	repo, err := gitLabVersioning.cloneRepoToDirectory()
	if err != nil {
		return nil, fmt.Errorf("error while initiating git package due to : %w", err)
	}

	gitLabVersioning.repo = repo

	return gitLabVersioning, nil
}
