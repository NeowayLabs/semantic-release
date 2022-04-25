package main

import (
	"flag"
	"os"

	"github.com/NeowayLabs/semantic-release/cmd/handler"
	"github.com/NeowayLabs/semantic-release/src/version"
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
		versionService := version.NewService(nil, nil)
		handlerService := handler.New(upgradeVersionCmd, gitHost, groupName, projectName, authKey, upgradePyFile, versionService)
		handlerService.HandleSemantic()
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
