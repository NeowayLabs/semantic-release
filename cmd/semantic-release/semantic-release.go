package main

import (
	"flag"
	"fmt"
	"os"

	commitmessage "github.com/NeowayLabs/semantic-release/src/commit-message"
	committype "github.com/NeowayLabs/semantic-release/src/commit-type"
	"github.com/NeowayLabs/semantic-release/src/files"
	"github.com/NeowayLabs/semantic-release/src/git"
	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/semantic"
	"github.com/NeowayLabs/semantic-release/src/time"
	v "github.com/NeowayLabs/semantic-release/src/version"
)

const (
	serviceName = "semantic-release"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

var (
	// version is set at build time
	Version  = "No version provided at build time"
	homePath = os.Getenv("HOME")
)

func main() {
	version := false
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()

	upgradeVersionCmd := flag.NewFlagSet("up", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpCommitCmd := flag.NewFlagSet("help-cmt", flag.ExitOnError)

	commitLint := upgradeVersionCmd.Bool("commit-lint", false, "Only lint commit history if set as true. (default false)")
	branchName := upgradeVersionCmd.String("branch-name", "", "Branch name to be cloned.")
	gitHost := upgradeVersionCmd.String("git-host", "", "Git host name. I.e.: gitlab.integration-tests.com. (required)")
	groupName := upgradeVersionCmd.String("git-group", "", "Git group name. (required)")
	projectName := upgradeVersionCmd.String("git-project", "", "Git project name. (required)")
	upgradePyFile := upgradeVersionCmd.Bool("setup-py", false, "Upgrade version in setup.py file. (default false)")
	username := upgradeVersionCmd.String("username", "", "Git username. (required)")
	password := upgradeVersionCmd.String("password", "", "Git password. (required)")
	logLevel := upgradeVersionCmd.String("log-level", "debug", "Log level.")

	if len(os.Args) < 2 {
		printWelcomeMessage()
		fmt.Println("\n" + colorRed + "Oops! Invalid input parameter." + colorCyan + " *** Usage: docker run neowaylabs/semantic-release [up] [help] [help-cmt] ***" + colorReset)
		os.Exit(1)
	}

	upgradeVersionCmd.Parse(os.Args[2:])

	if Version == "No version provided at build time" {
		Version = ""
	}

	logger, err := log.New(serviceName, Version, *logLevel)
	if err != nil {
		os.Exit(1)
	}

	printWelcomeMessage()
	switch os.Args[1] {
	case "up":
		logger.Info(colorYellow + "\nSemantic Version just started the process...\n\n" + colorReset)

		semantic := newSemantic(logger, upgradeVersionCmd, gitHost, groupName, projectName, username, password, upgradePyFile, branchName)

		if *commitLint {
			if *branchName == "" {
				logger.Error(colorRed + "\nThe argument -branch-name must be set when --commit-lint is true.\n\n" + colorReset)
			}

			logger.Info(colorYellow + "\nSemantic Version commit lint started...\n\n" + colorReset)
			err := semantic.CommitLint()
			if err != nil {
				printCommitTypes()
				printCommitMessageExample()
				os.Exit(1)
			}
		} else {
			if err := semantic.GenerateNewRelease(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
		}

		logger.Info(colorYellow + "\nDone!" + colorReset)
	case "help":
		printMainCommands()
		helpCmd.PrintDefaults()
		upgradeVersionCmd.PrintDefaults()

	case "help-cmt":
		helpCommitCmd.PrintDefaults()
		printCommitMessageExample()
		printCommitTypes()

	default:
		fmt.Printf(colorRed+"\nOops! Invalid input parameter [%v]. Expected [up], [help] or [help-cmt]."+colorReset+" \nRun "+colorCyan+"[docker run neowaylabs/semantic-release help`] "+colorReset+"to learn more about this CLI usage.\n", os.Args[1])
		os.Exit(1)
	}
}

type UpgradeFiles struct {
	Files []UpgradeFile
}

type UpgradeFile struct {
	Path            string
	DestinationPath string
	VariableName    string
}

func addFilesToUpgradeList(upgradePyFile *bool, repositoryRootPath string) UpgradeFiles {
	upgradeFilesList := UpgradeFiles{}
	if *upgradePyFile {
		upgradeFilesList.Files = append(upgradeFilesList.Files, UpgradeFile{Path: fmt.Sprintf("%s/setup.py", repositoryRootPath), DestinationPath: "", VariableName: "__version__"})
	}

	return upgradeFilesList
}

func validateIncomingParams(logger *log.Log, upgradeVersionCmd *flag.FlagSet, gitHost, groupName, projectName, username, password *string, upgradePyFile *bool) {
	if *gitHost == "" {
		logger.Info(colorRed + "Oops! Git host name must be specified." + colorReset + "[docker run neowaylabs/semantic-release up " + colorYellow + "-git-host gitHostNameHere]" + colorReset)
		os.Exit(1)
	}

	if *groupName == "" {
		logger.Info(colorRed + "Oops! Git group name must be specified." + colorReset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere " + colorYellow + "-git-group gitGroupNameHere]" + colorReset)
		os.Exit(1)
	}

	if *projectName == "" {
		logger.Info(colorRed + "Oops! Git project name must be specified." + colorReset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -git-group gitGroupNameHere " + colorYellow + "-git-project gitProjectNameHere]" + colorReset)
		os.Exit(1)
	}

	if *username == "" {
		logger.Info(colorRed + "Oops! Username must be specified." + colorReset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -git-group gitGroupNameHere -git-project gitProjectNameHere " + colorYellow + "-username gitUsername]" + colorReset)
		os.Exit(1)
	}

	if *password == "" {
		logger.Info(colorRed + "Oops! password must be specified." + colorReset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -git-group gitGroupNameHere -git-project gitProjectNameHere -username gitUsername " + colorYellow + "-password gitPassword]" + colorReset)
		os.Exit(1)
	}
}

func printWelcomeMessage() {
	fmt.Println(colorYellow + "\nWelcome to the Semantic Release CLI!" + colorReset)
	fmt.Println("\n\tThis CLI allows you to automatically upgrade a git project. \n\t\t* It changes the CHANGELOG.md file.\n\t\t* It Changes setup.py file (if setup-py parameter is set as true).\n\t\t* It also pushes the changes to master, creating and pushing a new corresponding tag.")
}

func printMainCommands() {
	fmt.Println(colorYellow + "\n\nHow to use it?" + colorReset)
	fmt.Println("\nThere are three main commands as follows:")
	fmt.Println(colorYellow + "\n\t* [docker run neowaylabs/semantic-release help]" + colorReset + ": this command shows you how to properly use the Semantic Release CLI.")
	fmt.Println(colorYellow + "\n\t* [docker run neowaylabs/semantic-release help-cmt]" + colorReset + ": this command shows you the commit types considered by the Semantic Release CLI.")
	fmt.Println(colorYellow + "\n\t* [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere -project gitProjectNameHere -username gitUsername -password gitPassword]" + colorReset + ": this command aims to automatically upgrade the project release version based on current commit subject.")
	fmt.Println("\nAvailable Parameters for " + colorYellow + "[docker run neowaylabs/semantic-release up]:" + colorReset)
}

func printCommitTypes() {
	fmt.Println(colorYellow + "\nTHE AVAILABLE COMMIT TYPES ARE:" + colorReset)
	fmt.Println(colorYellow + "\n\t*            [build]" + colorReset + ": Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)")
	fmt.Println(colorYellow + "\t*               [ci]" + colorReset + ": Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)")
	fmt.Println(colorYellow + "\t*             [docs]" + colorReset + ": Documentation only changes")
	fmt.Println(colorYellow + "\t*    [documentation]" + colorReset + ": ||")
	fmt.Println(colorYellow + "\t*             [feat]" + colorReset + ": A new feature")
	fmt.Println(colorYellow + "\t*          [feature]" + colorReset + ": ||")
	fmt.Println(colorYellow + "\t*              [fix]" + colorReset + ": A bug fix")
	fmt.Println(colorYellow + "\t*             [perf]" + colorReset + ": A code change that improves performance")
	fmt.Println(colorYellow + "\t*      [performance]" + colorReset + ": ||")
	fmt.Println(colorYellow + "\t*         [refactor]" + colorReset + ": A code change that neither fixes a bug nor adds a feature")
	fmt.Println(colorYellow + "\t*            [style]" + colorReset + ": Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)")
	fmt.Println(colorYellow + "\t*             [test]" + colorReset + ": Adding missing tests or correcting existing tests")
	fmt.Println(colorYellow + "\t*            [chore]" + colorReset + ": No production code change. It won't release a new version.")
	fmt.Println(colorYellow + "\t*             [skip]" + colorReset + ": Skip versioning. It won't release a new version.")
	fmt.Println(colorYellow + "\t*               [bc]" + colorReset + ": Changes that will require other changes in dependant applications")
	fmt.Println(colorYellow + "\t*         [breaking]" + colorReset + ": ||")
	fmt.Println(colorYellow + "\t*  [breaking change]" + colorReset + ": ||")
}

func printCommitMessageExample() {
	fmt.Println(colorYellow + "\nCOMMIT MESSAGE PATTERN" + colorReset)
	fmt.Println("\nThe commit message must follow the pattern below.")
	fmt.Println("\n\ttype(optional scope): Commit subject message here.")
	fmt.Println(colorYellow + "\n\tI.e." + colorReset)
	fmt.Println("\t\tfeat(config): Added new feature to handle configs.")

	fmt.Println("\n\tNote 1: The (scope) is optional. Semantic-release accepts the following pattern: \"type: Commit subject message here\".")
	fmt.Println("\n\tNote 2: The maximum number of characters is 150. If the commit subject exceeds it, it will be cut, keeping only the first 150 characters.")
}

func newSemantic(logger *log.Log, upgradeVersionCmd *flag.FlagSet, gitHost, groupName, projectName, username, password *string, upgradePyFile *bool, branchName *string) *semantic.Semantic {

	validateIncomingParams(logger, upgradeVersionCmd, gitHost, groupName, projectName, username, password, upgradePyFile)

	timer := time.New(logger)
	repositoryRootPath := fmt.Sprintf("%s/%s", homePath, *projectName)

	url := fmt.Sprintf("https://%s:%s@%s/%s/%s.git", *username, *password, *gitHost, *groupName, *projectName)
	repoVersionControl, err := git.New(logger, timer.PrintElapsedTime, url, *username, *password, repositoryRootPath, *branchName)
	if err != nil {
		logger.Fatal(err.Error())
	}

	commitTypeManager := committype.New(logger)
	commitMessageManager := commitmessage.New(logger, commitTypeManager)

	filesVersionControl := files.New(logger, timer.PrintElapsedTime, *gitHost, repositoryRootPath, *groupName, *projectName, commitMessageManager)

	versionControl := v.NewVersionControl(logger, timer.PrintElapsedTime, commitTypeManager)

	return semantic.New(logger, repositoryRootPath, addFilesToUpgradeList(upgradePyFile, repositoryRootPath), repoVersionControl, filesVersionControl, versionControl, commitMessageManager, commitTypeManager)
}
