package version

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	major          = "MAJOR"
	minor          = "MINOR"
	patch          = "PATCH"
	versionPattern = "%d.%d.%d"
	yellow         = "\033[33m"
	reset          = "\033[0m"

	ErrMsgConvertToInt = "could not convert %v to int"
)

var (
	commitChangeTypes             = []string{"build", "ci", "docs", "fix", "feat", "perf", "refactor", "style", "test", "breaking change", "breaking changes", "skip", "skip versioning", "skip v"}
	commitChangeTypesMajorUpgrade = []string{"breaking change", "breaking changes"}
	commitChangeTypesMinorUpgrade = []string{"feat"}
	commitChangeTypePatchUpgrade  = []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}
	commitTypeSkipVersioning      = []string{"skip", "skip versioning", "skip v"}
)

type ElapsedTime func(functionName string) func()

type VersionControl struct {
	printElapsedTime ElapsedTime
}

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
func (v *VersionControl) splitVersionMajorMinorPatch(version string) (map[string]int, error) {
	versionMap := make(map[string]string)
	resultMap := make(map[string]int)
	splitedVersion := strings.Split(version, ".")

	versionMap[major] = splitedVersion[0]
	versionMap[minor] = splitedVersion[1]
	versionMap[patch] = splitedVersion[2]

	for key, version := range versionMap {
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			log.Printf(ErrMsgConvertToInt, version)
			return nil, errors.New(fmt.Sprintf(ErrMsgConvertToInt, version))

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
func (v *VersionControl) getUpgradeType(commitChangeType string) (string, error) {
	if v.hasStringInSlice(commitChangeType, commitChangeTypesMajorUpgrade) {
		return major, nil
	} else if v.hasStringInSlice(commitChangeType, commitChangeTypesMinorUpgrade) {
		return minor, nil
	} else if v.hasStringInSlice(commitChangeType, commitChangeTypePatchUpgrade) {
		return patch, nil
	}
	return "", errors.New(fmt.Sprintf("%s is an invalid upgrade change type", commitChangeType))
}

// upgradeVersion upgrade the current version based on the upgradeType.
// Args:
// 		upgradeType (string): MAJOR, MINOR or PATCH.
// 		currentMajor (string): Current release major version. I.e.: >2<.1.1.
// 		currentMinor (string): Current release minor version. I.e.: 2.>1<.1.
// 		currentPatch (string): Current release patch version. I.e.: 2.1.>1<.
// Returns:
// 		It will return a string with the new version.
// 		I.e.:
// 		1 - If the current version is 2.1.1 and the update type is MAJOR it will return 3.0.0
// 		2 - If the current version is 2.1.1 and the update type is MINOR it will return 2.2.0
// 		1 - If the current version is 2.1.1 and the update type is PATCH it will return 2.1.2
func (v *VersionControl) upgradeVersion(upgradeType string, currentMajor, currentMinor, currentPatch int) string {
	var newVersion string

	switch upgradeType {
	case major:
		log.Printf(yellow+"%d"+reset+".0.0", currentMajor+1)
		newVersion = fmt.Sprintf(versionPattern, currentMajor+1, 0, 0)
	case minor:
		log.Printf("%d."+yellow+"%d"+reset+".0", currentMajor, currentMinor+1)
		newVersion = fmt.Sprintf(versionPattern, currentMajor, currentMinor+1, 0)
	case patch:
		log.Printf("%d.%d."+yellow+"%d"+reset, currentMajor, currentMinor, currentPatch+1)
		newVersion = fmt.Sprintf(versionPattern, currentMajor, currentMinor, currentPatch+1)
	}
	return newVersion
}

// GetNewVersion upgrade the current version based on the commitChangeType.
// It calls the getUpgradeType function to define where to upgrade the version (MAJOR.MINOR.PATCH).
// Args:
// 		commitMessage (string): The commit message.
// 		currentVersion (string): Current release version. I.e.: 2.1.1.
// Returns:
// 		It will return a string with the new version.
// 		I.e.:
// 		1 - If the current version is 2.1.1 and the update type is MAJOR it will return 3.0.0
// 		2 - If the current version is 2.1.1 and the update type is MINOR it will return 2.2.0
// 		1 - If the current version is 2.1.1 and the update type is PATCH it will return 2.1.2
func (v *VersionControl) GetNewVersion(commitMessage string, currentVersion string) (string, error) {
	defer v.printElapsedTime("GetNewVersion")()
	log.Printf("generating new version from %s", currentVersion)

	commitChangeType, err := v.getCommitChangeTypeFromMessage(commitMessage)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error while finding commit change type within commit message due to: %s", err.Error()))
	}

	curVersion, err := v.splitVersionMajorMinorPatch(currentVersion)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error while spliting version into MAJOR.MINOR.PATCH due to: %s", err.Error()))
	}
	currentMajor := curVersion[major]
	currentMinor := curVersion[minor]
	currentPatch := curVersion[patch]

	upgradeType, err := v.getUpgradeType(commitChangeType)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error while getting upgrade type due to: %s", err.Error()))
	}

	newVersion := v.upgradeVersion(upgradeType, currentMajor, currentMinor, currentPatch)

	return newVersion, nil
}

// getCommitChangeTypeFromMessage get the commit type from Message
// I.e.:
//       type: [fix]
//       message: Commit subject here.
// Output: fix
func (v *VersionControl) getCommitChangeTypeFromMessage(commitMessage string) (string, error) {
	log.Println("getting commit type from message")
	splitedMessage := strings.Split(commitMessage, "\n")
	for _, row := range splitedMessage {
		for _, changeType := range commitChangeTypes {
			if strings.Contains(strings.ToLower(row), "type:") && strings.Contains(strings.ToLower(row), fmt.Sprintf("[%s]", changeType)) {
				return changeType, nil
			}
		}
	}

	return "", errors.New("change type not found")
}

// hasStringInSlice aims to verify if a string is inside a slice of strings.
// It requires a full match.
// Args:
// 		value (string): String value to find.
// 		slice ([]string): Slice containing strings.
// Returns:
// 		bool: True when found, otherwise false.
func (v *VersionControl) hasStringInSlice(value string, slice []string) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}

// MustSkip compare commit type with skip types (CommitTypeSkipVersioning) to avoid upgrading version.
// I.e.:
//       commitChangeType: [skip]
// Output: true
func (v *VersionControl) MustSkipVersioning(commitMessage string) bool {
	commitChangeType, err := v.getCommitChangeTypeFromMessage(commitMessage)
	if err != nil {
		return true
	}

	return v.hasStringInSlice(commitChangeType, commitTypeSkipVersioning)
}

// NewVersionControl is the version control constructor
func NewVersionControl(printElapsedTime ElapsedTime) *VersionControl {
	return &VersionControl{
		printElapsedTime: printElapsedTime,
	}
}
