package consts

import (
	"fmt"

	style "github.com/NeowayLabs/semantic-release/src/utils/stdout_style"
)

const (
	Major                         = "MAJOR"
	Minor                         = "MINOR"
	Patch                         = "PATCH"
	VersionPattern                = "%d.%d.%d"
	VersionPyVariable             = "__version__"
	ChangeLogDefaultFile          = "CHANGELOG.md"
	SetupPythonDefaultFile        = "setup.py"
	CommitMsgLimitCharacters      = 150
	AuthorChangelogFormat         = "@%s"
	ChangeLogCommitHashLinkFormat = "[%s](https://%s/%s/%s/-/commit/%s)"
)

var (
	CommitChangeTypes             = []string{"build", "ci", "docs", "fix", "feat", "perf", "refactor", "style", "test", "breaking change", "skip", "skip versioning", "skip v"}
	CommitChangeTypesMajorUpgrade = []string{"breaking change"}
	CommitChangeTypesMinorUpgrade = []string{"feat"}
	CommitChangeTypePatchUpgrade  = []string{"build", "ci", "docs", "fix", "perf", "refactor", "style", "test"}
	CommitTypeSkipVersioning      = []string{"skip", "skip versioning", "skip v"}
)

func PrintCommitTypes() {
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
	fmt.Println(style.Yellow + "\t*          [skip v]" + style.Reset + ": Skip versioning")
	fmt.Println(style.Yellow + "\t* [skip versioning]" + style.Reset + ": Skip versioning")
}
