package git_test

import (
	"fmt"

	"github.com/NeowayLabs/semantic-release/src/git"
	"github.com/NeowayLabs/semantic-release/src/log"
)

type GitLabVersioningMock struct {
	url                  string
	destinationDirectory string
	sshKey               string
	replace              bool
	username             string
	password             string
}

func printElapsedTimeMock(functionName string) func() {
	return func() {
		fmt.Printf("%s done.", functionName)
	}
}

type fixture struct {
	gitLabVersioning *GitLabVersioningMock
	log              *log.Log
}

func setup() *fixture {
	logger, err := log.New("test", "1.0.0", "debug")
	if err != nil {
		panic(err.Error())
	}
	return &fixture{log: logger, gitLabVersioning: &GitLabVersioningMock{}}
}

func (f *fixture) NewGitService() (*git.GitVersioning, error) {
	return git.New(f.log, printElapsedTimeMock, f.gitLabVersioning.url, f.gitLabVersioning.username, f.gitLabVersioning.password, f.gitLabVersioning.destinationDirectory)
}
