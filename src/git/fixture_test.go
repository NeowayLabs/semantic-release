package git_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/git"
	"github.com/NeowayLabs/semantic-release/src/log"
)

type GitLabVersioningMock struct {
	url                  string
	destinationDirectory string
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

func getValidSetup() *fixture {
	f := setup()
	f.gitLabVersioning.url = "https://gitlab/dataplatform/integration-tests.git"
	f.gitLabVersioning.destinationDirectory = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "integration-tests")
	f.gitLabVersioning.username = "root"
	f.gitLabVersioning.password = "password"
	return f
}

func newValidVersion(t *testing.T, currentVersion string) string {
	splitedVersion := strings.Split(currentVersion, ".")
	var newVersionSlice []int
	for _, version := range splitedVersion {
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			t.Errorf("could not convert %v to int", version)
			return ""
		}
		newVersionSlice = append(newVersionSlice, versionInt)
	}

	newVersionSlice[0] = newVersionSlice[0] + 1
	return fmt.Sprintf("%v.%v.%v", newVersionSlice[0], newVersionSlice[1], newVersionSlice[2])
}

func (f *fixture) isDangerousOsOperation() bool {
	if strings.ReplaceAll(strings.ReplaceAll(f.gitLabVersioning.destinationDirectory, "/", ""), " ", "") == strings.ReplaceAll(strings.ReplaceAll(os.Getenv("HOME"), "/", ""), " ", "") {
		return true
	}
	return false
}

func (f *fixture) cleanLocalRepo(t *testing.T) {
	if f.isDangerousOsOperation() {
		t.Errorf("error while cleaning up local repository path %s due to: too danger os operation", f.gitLabVersioning.destinationDirectory)
	}

	if _, err := os.Stat(f.gitLabVersioning.destinationDirectory); os.IsNotExist(err) {
		t.Errorf("error while cleaning removing repository path %s due to: %s", f.gitLabVersioning.destinationDirectory, err.Error())
	}

	err := os.RemoveAll(f.gitLabVersioning.destinationDirectory)
	if err != nil {
		t.Errorf("error while cleaning removing repository path %s due to: %s", f.gitLabVersioning.destinationDirectory, err.Error())
	}
}

func (f *fixture) newGitService() (*git.GitVersioning, error) {
	return git.New(f.log, printElapsedTimeMock, f.gitLabVersioning.url, f.gitLabVersioning.username, f.gitLabVersioning.password, f.gitLabVersioning.destinationDirectory)
}
