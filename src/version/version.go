package version

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	e "github.com/NeowayLabs/semantic-release/src/errors"
	"github.com/NeowayLabs/semantic-release/src/slice"
	style "github.com/NeowayLabs/semantic-release/src/style"
	"github.com/NeowayLabs/semantic-release/src/time"
)

const (
	major                    = "MAJOR"
	minor                    = "MINOR"
	patch                    = "PATCH"
	versionPattern           = "%d.%d.%d"
	versionPyVariable        = "__version__"
	commitMsgLimitCharacters = 150
	authorChangelogFormat    = "@%s"
)

var (
	commitChangeTypes             = []string{"build", "ci", "docs", "fix", "feat", "perf", "refactor", "style", "test", "breaking change", "breaking changes", "skip", "skip versioning", "skip v"}
	commitChangeTypesMajorUpgrade = []string{"breaking change", "breaking changes"}
	commitChangeTypesMinorUpgrade = []string{"feat"}
	commitChangeTypePatchUpgrade  = []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}
	commitTypeSkipVersioning      = []string{"skip", "skip versioning", "skip v"}
)

type Version interface {
	GetNewReleaseVersion(currentVersion string, commitChangeType string) (string, error)
	GetUpgradeType(commitChangeType string) (string, error)
	PrettifyCommitMessage(commitMessage string) (string, error)
	PrettifyAuthorEmailForChangelog(authorEmail string) string
	GetCommitChangeTypeFromMessage(commitMessage string) (string, error)
	MustSkipVersioning(commitChangeType string) bool
}

type version struct{}

// splitVersionMajorMinorPatch get a string version, split it and return a map of int values
// Args:
// 		version (string): Version to be splited. I.e: 2.1.1
// Returns:
// 		Success:
// 		It returns a map of int values
//      I.e.: map[MAJOR:2 MINOR:1 PATCH:1]
//
// 		Otherwise:
// 			error
func splitVersionMajorMinorPatch(version string) (map[string]int, error) {
	versionMap := make(map[string]string)
	resultMap := make(map[string]int)
	splitedVersion := strings.Split(version, ".")

	versionMap[major] = splitedVersion[0]
	versionMap[minor] = splitedVersion[1]
	versionMap[patch] = splitedVersion[2]

	for key, version := range versionMap {
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			log.Printf(e.ErrMsgConvertToInt, version)
			return nil, errors.New(fmt.Sprintf(e.ErrMsgConvertToInt, version))

		}
		resultMap[key] = versionInt
	}

	return resultMap, nil
}

// getUpgradeType defines where to update the current version
// MAJOR.MINOR.PATCH. I.e: 2.1.1
// Args:
// 		commitChangeType (string): Type of changes within the commit. I.e.: fix, feat, doc, etc. Take a look at CommitChangeTypes variable.
// Returns:
// 		MAJOR: if the commit type is in CommitChangeTypesMajorUpgrade slice
// 		MINOR: if the commit type is in CommitChangeTypesMinorUpgrade slice
// 		PATCH: if the commit type is in CommitChangeTypePatchUpgrade slice
// 		Otherwise, it returns an error
func (s *version) GetUpgradeType(commitChangeType string) (string, error) {
	if slice.IsStringInSlice(commitChangeType, commitChangeTypesMajorUpgrade) {
		return major, nil
	} else if slice.IsStringInSlice(commitChangeType, commitChangeTypesMinorUpgrade) {
		return minor, nil
	} else if slice.IsStringInSlice(commitChangeType, commitChangeTypePatchUpgrade) {
		return patch, nil
	}
	return "", errors.New(fmt.Sprintf(e.ErrMsgInvalidCommitChangeType, commitChangeType))
}

// getNewReleaseVersion upgrade the current version based on the commitChangeType.
// It calls the getUpgradeType function to define where to upgrade the version (MAJOR.MINOR.PATCH).
// Args:
// 		currentVersion (string): Current release version. I.e.: 2.1.1.
// 		commitChangeType (string): Type of changes within the commit. I.e.: fix, feat, doc, etc. Take a look at CommitChangeTypes variable.
// Returns:
// 		It will return a string with the new version.
// 		I.e.:
// 		1 - If the current version is 2.1.1 and the update type is MAJOR it will return 3.0.0
// 		2 - If the current version is 2.1.1 and the update type is MINOR it will return 2.2.0
// 		1 - If the current version is 2.1.1 and the update type is PATCH it will return 2.1.2
func (s *version) GetNewReleaseVersion(currentVersion string, commitChangeType string) (string, error) {
	log.Printf("generating new version from %s", currentVersion)
	var newVersion string
	curVersion, errSplitVersion := splitVersionMajorMinorPatch(currentVersion)
	if errSplitVersion != nil {
		return "", errors.New(fmt.Sprintf(e.ErrMsgSplitVersion, errSplitVersion))
	}
	currentMajor := curVersion[major]
	currentMinor := curVersion[minor]
	currentPatch := curVersion[patch]

	upgradeType, errUpgradeType := s.GetUpgradeType(commitChangeType)
	if errUpgradeType != nil {
		return "", errors.New(fmt.Sprintf(e.ErrMsgGetUpgradeType, errUpgradeType))
	}
	log.Printf("Commit type "+style.Green+"%s "+style.Reset+" requires a "+style.Yellow+"%s "+style.Reset+"version update", commitChangeType, upgradeType)
	switch upgradeType {
	case major:
		log.Printf(style.Yellow+"%d"+".0.0", currentMajor+1)
		newVersion = fmt.Sprintf(versionPattern, currentMajor+1, 0, 0)
	case minor:
		log.Printf("%d."+style.Yellow+"%d"+style.Reset+".0", currentMajor, currentMinor+1)
		newVersion = fmt.Sprintf(versionPattern, currentMajor, currentMinor+1, 0)
	case patch:
		log.Printf("%d.%d."+style.Yellow+"%d"+style.Reset, currentMajor, currentMinor)
		newVersion = fmt.Sprintf(versionPattern, currentMajor, currentMinor, currentPatch+1)
	default:
		newVersion = ""
	}
	if newVersion == "" {
		return "", errors.New(e.ErrMsgNewReleaseVersion)
	}
	return newVersion, nil
}

// PrettifyCommitMessage aims to keep a short message based on the commit message.
// Args:
// 		commitMessage (string): Full commit message.
// Returns:
// 		string: Returns a commit message with limmited number of characters.
// 		err: Error whenever unexpected issues happen.
func (s *version) PrettifyCommitMessage(commitMessage string) (string, error) {
	messageTag := "message:"
	if !strings.Contains(strings.ToLower(commitMessage), messageTag) {
		return "", errors.New(e.ErrMsgGetCommitMsgPatternNotFound)
	}

	var message string
	splitedMessage := strings.Split(commitMessage, "\n")

	for _, row := range splitedMessage {
		row := strings.ToLower(row)
		if strings.Contains(row, messageTag) {
			messagePosition := strings.Index(row, messageTag) + len(messageTag)
			message = strings.TrimSpace(row[messagePosition:])
		}
	}

	if len(message) >= 150 {
		message = fmt.Sprintf("%s...", message[:commitMsgLimitCharacters])
	}

	message = fmt.Sprintf("%s%s", strings.ToUpper(message[:1]), message[1:])
	return message, nil
}

// PrettifyAuthorEmailForChangelog gets the author email and returns in changelog format
// Args:
// 		authorEmail (string): Commit author email.
// Returns:
// 		string: return a string with the author email formated according to AuthorChangelogFormat
// 		I.e.:
//  		PrettifyAuthorEmailForChangelog("user.name@neoway.com.br")
// 			Result: @user.name
func (s *version) PrettifyAuthorEmailForChangelog(authorEmail string) string {
	atPosition := strings.Index(authorEmail, "@")
	return fmt.Sprintf(authorChangelogFormat, authorEmail[:atPosition])
}

// GetCommitChangeTypeFromMessage get the commit type from Message
// I.e.:
//       type: [fix]
//       message: Commit subject here.
// Output: fix
func (s *version) GetCommitChangeTypeFromMessage(commitMessage string) (string, error) {
	defer time.GetElapsedTime("GetCommitChangeTypeFromMessage")()
	log.Println("getting commit type from message")
	splitedMessage := strings.Split(commitMessage, "\n")
	for _, row := range splitedMessage {
		for _, changeType := range commitChangeTypes {
			if strings.Contains(strings.ToLower(row), "type:") && strings.Contains(strings.ToLower(row), fmt.Sprintf("[%s]", changeType)) {
				return changeType, nil
			}
		}
	}

	return "", errors.New(e.ErrMsgGetChangeType)
}

// MustSkip compare commit type with skip types (CommitTypeSkipVersioning) to avoid upgrading version.
// I.e.:
//       commitChangeType: [skip]
// Output: true
func (s *version) MustSkipVersioning(commitChangeType string) bool {
	return slice.IsStringInSlice(commitChangeType, commitTypeSkipVersioning)
}

// NewService is the version service constructor
func NewService() Version {
	return &version{}
}
