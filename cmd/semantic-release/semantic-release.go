package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/NeowayLabs/semantic-release/src/handler"
	style "github.com/NeowayLabs/semantic-release/src/style"
)

func main() {

	upgradeVersionCmd := flag.NewFlagSet("up", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpCommitCmd := flag.NewFlagSet("help-cmt", flag.ExitOnError)

	gitHost := upgradeVersionCmd.String("git-host", "", "Git host name. I.e.: gitlab.integration-tests.com. (required)")
	groupName := upgradeVersionCmd.String("git-group", "", "Git group name. (required)")
	projectName := upgradeVersionCmd.String("git-project", "", "Git project name. (required)")
	upgradePyFile := upgradeVersionCmd.Bool("setup-py", false, "Upgrade version in setup.py file. (default false)")
	upgradeChangelog := upgradeVersionCmd.Bool("changelog", true, "Upgrade version in CHANGELOG.md file.")
	createAndPushTag := upgradeVersionCmd.Bool("create-git-tag", true, "Create and push a new git tag.")
	pushChanges := upgradeVersionCmd.Bool("push", true, "Push changed files to the branch.")
	authKey := upgradeVersionCmd.String("auth", "", "SSH key. (required)")

	if len(os.Args) < 2 {
		printWelcomeMessage()
		log.Println("\n" + style.Red + "Oops! Invalid input parameter." + style.Cyan + " *** Usage: docker run neowaylabs/semantic-release [up] [help] [help-cmt] ***" + style.Reset)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		fmt.Println(style.Yellow + "\nSemantic Version just started the process...\n\n" + style.Reset)
		fmt.Println(authKey)
		fmt.Println(upgradePyFile)
		handler.HandleVersioning(upgradeVersionCmd, gitHost, groupName, projectName, authKey, upgradePyFile, upgradeChangelog, createAndPushTag, pushChanges)
	case "help":
		printWelcomeMessage()
		helpCmd.PrintDefaults()
		fmt.Println(style.Yellow + "\n\nHow to use it?" + style.Reset)
		fmt.Println("\nThere are three main commands as follows:")
		fmt.Println(style.Yellow + "\n\t* [docker run neowaylabs/semantic-release help]" + style.Reset + ": this command shows you how to properly use the Semantic Release CLI.")
		fmt.Println(style.Yellow + "\n\t* [docker run neowaylabs/semantic-release help-cmt]" + style.Reset + ": this command shows you the commit types considered by the Semantic Release CLI.")
		fmt.Println(style.Yellow + "\n\t* [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere -project gitProjectNameHere]" + style.Reset + ": this command aims to automatically upgrade the project release version based on current commit subject.")
		fmt.Println("\nAvailable Parameters for " + style.Yellow + "[docker run neowaylabs/semantic-release up]:" + style.Reset)
		upgradeVersionCmd.PrintDefaults()
	case "help-cmt":
		printWelcomeMessage()
		helpCommitCmd.PrintDefaults()
		printCommitTypes()
	default:
		printWelcomeMessage()
		fmt.Printf(style.Red+"\nOops! Invalid input parameter [%v]. Expected [up], [help] or [help-cmt]."+style.Reset+" \nRun "+style.Cyan+"[docker run neowaylabs/semantic-release help`] "+style.Reset+"to learn more about this CLI usage.\n", os.Args[1])
		os.Exit(1)
	}
}

func printWelcomeMessage() {
	fmt.Println(style.Yellow + "\nWelcome to the Semantic Release CLI!" + style.Reset)
	fmt.Println("\n\tThis CLI allows you to automatically upgrade a git project. \n\t\t* It changes the CHANGELOG.md file.\n\t\t* It Changes setup.py file (if setup-py parameter is set as true).\n\t\t* It also pushes the changes to master, create and pushing a new tag to master.")
}

func printCommitTypes() {
	fmt.Println(style.Yellow + "\nTHE AVAILABLE COMMIT TYPES ARE:" + style.Reset)
	fmt.Println(style.Yellow + "\n\t*           [build]" + style.Reset + ": Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)")
	fmt.Println(style.Yellow + "\t*              [ci]" + style.Reset + ": Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)")
	fmt.Println(style.Yellow + "\t*            [docs]" + style.Reset + ": Documentation only changes")
	fmt.Println(style.Yellow + "\t*            [feat]" + style.Reset + ": A new feature")
	fmt.Println(style.Yellow + "\t*             [fix]" + style.Reset + ": A bug fix")
	fmt.Println(style.Yellow + "\t*            [perf]" + style.Reset + ": A code change that improves performance")
	fmt.Println(style.Yellow + "\t*        [refactor]" + style.Reset + ": A code change that neither fixes a bug nor adds a feature")
	fmt.Println(style.Yellow + "\t*           [style]" + style.Reset + ": Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)")
	fmt.Println(style.Yellow + "\t*            [test]" + style.Reset + ": Adding missing tests or correcting existing tests")
	fmt.Println(style.Yellow + "\t*            [skip]" + style.Reset + ": Skip versioning")
	fmt.Println(style.Yellow + "\t* [skip versioning]" + style.Reset + ": Skip versioning")
	fmt.Println(style.Yellow + "\t* [breaking change]" + style.Reset + ": Change that will require other changes in dependant applications")
	fmt.Println(style.Yellow + "\t* [breaking changes]" + style.Reset + ": Changes that will require other changes in dependant applications")
}
