package errorsutils

import (
	"fmt"
	"log"
)

const (
	ErrMsgCloneFail                      = "error while clonning repository"
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
	ErrMsgGettingMostRecentCommit        = "error while getting most recent commit"
	ErrMsgRepoAlreadyCloned              = "repository was already cloned"
	ErrMsgDueTo                          = "%s due to: %s"
	ErrMsgUpgradeSetupPython             = "error while upgrading setup.py file due to: %s"
	ErrMsgUpgradeChangelog               = "error while upgrading CHANGELOG.md file due to: %s"
	ErrMsgNoSuchFileOrDirectory          = "no such file or directory"
	ErrMsgWritingUpgradedFile            = "\n\nerror while writing %s file with new version %s due to: %s"
	ErrMsgScanningFile                   = "\n\nerror while scanning file: %s due to: %s"
	ErrMsgNewReleaseVersion              = "error while getting new release version"
	ErrMsgConvertToInt                   = "could not convert %v to int"
	ErrMsgSplitVersion                   = "error while spliting version into MAJOR.MINOR.PATCH due to: %s"
	ErrMsgGetUpgradeType                 = "error while getting upgrade type due to: %s"
)

func Error(err error, message string) bool {
	if err != nil {
		log.Printf(fmt.Sprintf(ErrMsgDueTo, message, err))
		return true
	}
	return false
}

func HasError(err error, message string) {
	if Error(err, message) {
		panic(err)
	}
}
