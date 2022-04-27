package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/NeowayLabs/semantic-release/src/semantic"
)

var (
	homePath = os.Getenv("HOME")
)

func main() {
	upgradeVersionCmd := flag.NewFlagSet("up", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpCommitCmd := flag.NewFlagSet("help-cmt", flag.ExitOnError)

	gitHost := upgradeVersionCmd.String("git-host", "", "Git host name. I.e.: gitlab.integration-tests.com. (required)")
	groupName := upgradeVersionCmd.String("git-group", "", "Git group name. (required)")
	projectName := upgradeVersionCmd.String("git-project", "", "Git project name. (required)")
	upgradePyFile := upgradeVersionCmd.Bool("setup-py", false, "Upgrade version in setup.py file. (default false)")
	authKey := upgradeVersionCmd.String("auth", "", "SSH key. (required)")

	if len(os.Args) < 2 {
		// TODO: Implement me!
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		// TODO: Implement me!
		// TODO: Implement me!
		upgradeVersionCmd.Parse(os.Args[2:])

		validate(upgradeVersionCmd, gitHost, groupName, projectName, authKey, upgradePyFile)

		versionService := semantic.New(homePath, addFilesToUpgradeList(upgradePyFile), nil, nil, nil)
		fmt.Println(versionService)
		// TODO: Implement me!
	case "help":
		// TODO: Implement me!
		helpCmd.PrintDefaults()
		upgradeVersionCmd.PrintDefaults()
	case "help-cmt":
		// TODO: Implement me!
		helpCommitCmd.PrintDefaults()

	default:
		// TODO: Implement me!
		os.Exit(1)
	}
}

type upgradeFiles struct {
	files []upgradeFile
}

type upgradeFile struct {
	path         string
	variableName string
}

func addFilesToUpgradeList(upgradePyFile *bool) upgradeFiles {
	upgradeFilesList := upgradeFiles{}
	if *upgradePyFile {
		upgradeFilesList.files = append(upgradeFilesList.files, upgradeFile{path: homePath + "setup.py", variableName: "__version__"})
	}
	return upgradeFilesList
}

func validate(upgradeVersionCmd *flag.FlagSet, gitHost, groupName, projectName, authKey *string, upgradePyFile *bool) {

	if *gitHost == "" {
		log.Println("Oops! Git host name must be specified. [docker run neowaylabs/semantic-release up -git-host gitHostNameHere]")
		os.Exit(1)
	}

	if *groupName == "" {
		log.Println("Oops! Git group name must be specified. [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere]")
		os.Exit(1)
	}

	if *projectName == "" {
		log.Println("Oops! Git project name must be specified. [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere -project gitProjectNameHere]")
		os.Exit(1)
	}

	if *authKey == "" {
		log.Println("Oops! Auth must be specified. [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere -project gitProjectNameHere -auth sshKeyHere]")
		os.Exit(1)
	}
}
