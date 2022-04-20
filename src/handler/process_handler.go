package handler

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/NeowayLabs/semantic-release/src/auth"
	e "github.com/NeowayLabs/semantic-release/src/errors"
	"github.com/NeowayLabs/semantic-release/src/files"
	"github.com/NeowayLabs/semantic-release/src/git"
	style "github.com/NeowayLabs/semantic-release/src/style"
	"github.com/NeowayLabs/semantic-release/src/version"
)

const (
	homeEnv                = "HOME"
	changeLogDefaultFile   = "CHANGELOG.md"
	setupPythonDefaultFile = "setup.py"
)

func HandleVersioning(upgradeVersionCmd *flag.FlagSet, gitHost, groupName, projectName, authKey *string, upgradePyFile, upgradeChangelog, createAndPushTag, pushChanges *bool) {
	upgradeVersionCmd.Parse(os.Args[2:])

	if *gitHost == "" {
		log.Println(style.Red + "Oops! Git host name must be specified." + style.Reset + " [docker run neowaylabs/semantic-release up " + style.Yellow + "-git-host" + " gitHostNameHere]" + style.Reset)
		os.Exit(1)
	}

	if *groupName == "" {
		log.Println(style.Red + "Oops! Git group name must be specified." + style.Reset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere " + style.Yellow + "-group" + " gitGroupNameHere]" + style.Reset)
		os.Exit(1)
	}

	if *projectName == "" {
		log.Println(style.Red + "Oops! Git project name must be specified." + style.Reset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere " + style.Yellow + "-project" + " gitProjectNameHere]" + style.Reset)
		os.Exit(1)
	}

	fmt.Println(*upgradePyFile)
	if *authKey == "" {
		log.Println(style.Red + "Oops! Auth must be specified." + style.Reset + " [docker run neowaylabs/semantic-release up -git-host gitHostNameHere -group gitGroupNameHere -project gitProjectNameHere" + style.Yellow + " -auth" + " sshKeyHere]" + style.Reset)
		os.Exit(1)
	}

	homePath := os.Getenv(homeEnv)

	repoUrl := fmt.Sprintf("git@%s:%s/%s.git", *gitHost, *groupName, *projectName)
	gitService := git.New()

	sshKey, errPublicKey := auth.GetPublicKey([]byte(auth.FormatSSHKey(*authKey, "#")))
	e.HasError(errPublicKey, e.ErrMsgGettingPublicKey)

	destinationPath := fmt.Sprintf("%s/%s", homePath, *projectName)
	localRepo, errClone := gitService.CloneRepoToDirectory(repoUrl, destinationPath, true, sshKey)
	e.HasError(errClone, e.ErrMsgCloneFail)

	refBranchHead, errBranch := gitService.GetBranchPointedToHead(localRepo)
	e.HasError(errBranch, e.ErrMsgRetrievingBranchHead)

	commitHistory, errCmtHistory := gitService.GetCommitHistory(localRepo, refBranchHead)
	e.HasError(errCmtHistory, e.ErrMsgRetrievingCommitHistory)

	var commitList []interface{}
	for _, commit := range commitHistory {
		commitList = append(commitList, commit)
	}
	mostRecentCommit, errRecentCommit := gitService.GetMostRecentCommit(commitList)
	e.HasError(errRecentCommit, e.ErrMsgGettingMostRecentCommit)

	fmt.Printf(style.Yellow+"MOST RECENT COMMIT:"+style.Reset+"\n\t%s\n", mostRecentCommit)

	versionService := version.NewService()
	commitChangeType, errCommitChangeType := versionService.GetCommitChangeTypeFromMessage(mostRecentCommit.Message)
	e.HasError(errCommitChangeType, e.ErrMsgGetCommitChangeTypeFromMessage)
	fmt.Printf(style.Yellow+"CURRENT CHANGE TYPE "+style.Reset+">>> "+style.BGRed+"[%s]"+style.Reset+"\n\n", commitChangeType)

	if !versionService.MustSkipVersioning(commitChangeType) {
		mostRecentTag, errRecentTag := gitService.GetMostRecentTag(localRepo)
		e.HasError(errRecentTag, e.ErrMsgGettingMostRecentTag)

		newReleaseVersion, errNewReleaseVersion := versionService.GetNewReleaseVersion(mostRecentTag, commitChangeType)
		e.HasError(errNewReleaseVersion, e.ErrMsgNewReleaseVersion)
		fmt.Printf(style.Yellow+"NEW RELEASE VERSION "+style.Reset+">>> "+style.BGRed+"%s"+style.Reset+"\n\n", newReleaseVersion)

		// Update setup.py file
		if *upgradePyFile {
			file := files.File{
				OriginPath:        fmt.Sprintf("%s/%s/%s", homePath, *projectName, setupPythonDefaultFile),
				OutputPath:        fmt.Sprintf("%s/%s/%s", homePath, *projectName, setupPythonDefaultFile),
				NewReleaseVersion: newReleaseVersion,
			}
			filesService := files.New(file)

			errUpgradeSetupPython := filesService.UpgradeVersionInSetupPyFile()
			e.HasError(errUpgradeSetupPython, fmt.Sprintf(e.ErrMsgUpgradeSetupPython, errUpgradeSetupPython))
		}

		// Update CHANGELOG.md
		if *upgradeChangelog {
			commitMessage, errCommitMsg := versionService.PrettifyCommitMessage(mostRecentCommit.Message)
			e.HasError(errCommitMsg, e.ErrMsgGetCommitMsgPatternNotFound)

			file := files.File{
				OriginPath:        fmt.Sprintf("%s/%s/%s", homePath, *projectName, changeLogDefaultFile),
				OutputPath:        fmt.Sprintf("%s/%s/%s", homePath, *projectName, changeLogDefaultFile),
				NewReleaseVersion: newReleaseVersion,
			}
			filesService := files.New(file)
			errUpgradeChangeLog := filesService.UpgradeChangelogFile(*gitHost, commitChangeType, mostRecentCommit.Hash.String(), commitMessage, versionService.PrettifyAuthorEmailForChangelog(mostRecentCommit.Author.Email), *groupName, *projectName)
			e.HasError(errUpgradeChangeLog, e.ErrMsgUpgradeChangelog)
		}

		errAdd := gitService.AddToStage(localRepo)
		e.HasError(errAdd, e.ErrMsgAddingToStage)

		errCommit := gitService.CommitChanges(localRepo, mostRecentCommit.Author.Name, mostRecentCommit.Author.Email, fmt.Sprintf("type: [skip]: message: Commit automatically generated by Semantic Release. The new tag is %s", newReleaseVersion))
		e.HasError(errCommit, e.ErrMsgCommitingToTransfer)

		errPushCommit := gitService.Push(localRepo, sshKey)
		e.HasError(errPushCommit, e.ErrMsgPushCommits)

		// Create and push new git TAG
		if *createAndPushTag {
			tagCreated, errCreateTag := gitService.SetTag(localRepo, newReleaseVersion, mostRecentCommit.Author.Name, mostRecentCommit.Author.Email)
			e.HasError(errCreateTag, e.ErrMsgCreateTag)

			if tagCreated {
				errPush := gitService.PushTags(localRepo, sshKey)
				e.HasError(errPush, e.ErrMsgCreateTag)
			}
		}
		log.Println(style.Purple + "The entire process worked like a charm.\nYou can check that by accessing the default branch in your repository.\n\n\t\t THANK YOU!" + style.Reset)
	}
}
