package errors

import (
	"fmt"
	"log"
)

const (
	ErrMsgCloneFail                      = "error while cloning repository"
	ErrMsgRetrievingBranchHead           = "error while retrieving the branch pointed by HEAD"
	ErrMsgRetrievingCommitHistory        = "error while retrieving the commit history"
	ErrMsgIteratingCommitHistory         = "error while iterating over the commits"
	ErrMsgGetChangeType                  = "change type not found"
	ErrMsgInvalidCommitChangeType        = "%s is an invalid commit change type"
	ErrMsgRetrievingTags                 = "error while retrieving tags from repository"
	ErrMsgGettingMostRecentTag           = "error while getting most recent tag"
	ErrMsgGetCommitChangeTypeFromMessage = "error while finding commit change type within commit message"
	ErrMsgGetCommitMsgPatternNotFound    = "commit message has no tag 'message:'"
	ErrMsgNoCommitsFound                 = "no commits found"
	ErrMsgInterfaceToStruct              = "error while casting interface to struct"
	ErrMsgGettingPublicKey               = "error while getting public key"
	ErrMsgGettingMostRecentCommit        = "error while getting most recent commit"
	ErrMsgRepoAlreadyCloned              = "repository was already cloned"
	ErrMsgDueTo                          = "%s due to: %s"
	ErrMsgUpgradeSetupPython             = "error while upgrading setup.py file due to: %s"
	ErrMsgUpgradeChangelog               = "error while upgrading CHANGELOG.md file due to: %s"
	ErrMsgAddingToStage                  = "error while adding changes to stage"
	ErrMsgCommitingToTransfer            = "error while commiting changes to tranfer area"
	ErrMsgNoSuchFileOrDirectory          = "no such file or directory"
	ErrMsgWritingUpgradedFile            = "\n\nerror while writing %s file with new version %s due to: %s"
	ErrMsgScanningFile                   = "\n\nerror while scanning file: %s due to: %s"
	ErrMsgNewReleaseVersion              = "error while getting new release version"
	ErrMsgConvertToInt                   = "could not convert %v to int"
	ErrMsgSplitVersion                   = "error while spliting version into MAJOR.MINOR.PATCH due to: %s"
	ErrMsgGetUpgradeType                 = "error while getting upgrade type due to: %s"
	ErrMsgCreateTag                      = "error while creating tag"
	ErrMsgTagExists                      = "tag %s already exists"
	ErrMsgTagDoesNotExists               = "tag %s does not exists"
	ErrMsgPushTags                       = "error while pushing tags"
	ErrMsgPushCommits                    = "error while pushing commits"
	ErrMsgIterateTagsError               = "error while iterating tags due to: %s"
	ErrMsgGetingTags                     = "error while getting tags due to %s"
)

// Error abstract the error verification, add log when error and return a boolean
// Args:
// 		err (error): Error to be checked.
// 		message (string): Message to be logged when error.
// Returns:
// 		bool: true when error, otherwise false.
func Error(err error, message string) bool {
	if err != nil {
		log.Printf(fmt.Sprintf(ErrMsgDueTo, message, err))
		return true
	}
	return false
}

// HasError aims to handle the CLI errors, logging and panic
// Args:
// 		err (error): Error to be checked.
// 		message (string): Message to be logged when error.
func HasError(err error, message string) {
	if Error(err, message) {
		panic(err)
	}
}
