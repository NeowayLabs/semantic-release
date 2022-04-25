package version

type RepositoryVersionControl interface {
	GetMessage() (string, error)
	GetVersionUpdateType(message string) (string, error)
	UpgradeRemoteRepository() error
}

type FilesVersionControl interface {
	UpgradeChangeLog(path string) error
	UpgradeVariableInFiles(variableName string, filesList []string) error
}

type Version struct {
	repoVersionControl  RepositoryVersionControl
	filesVersionControl FilesVersionControl
}

func NewService(repoVersionControl RepositoryVersionControl, filesVersionControl FilesVersionControl) *Version {
	return &Version{
		repoVersionControl:  repoVersionControl,
		filesVersionControl: filesVersionControl,
	}
}
