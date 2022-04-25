package handler

import (
	"flag"

	"github.com/NeowayLabs/semantic-release/src/version"
)

type HandlerService interface {
	HandleSemantic()
}

type Handler struct {
	upgradeVersionCmd *flag.FlagSet
	gitHost           *string
	groupName         *string
	projectName       *string
	authKey           *string
	upgradePyFile     *bool
	versionService    *version.Version
}

func (h *Handler) HandleSemantic() {
	// TODO: Implement me!
}

func New(upgradeVersionCmd *flag.FlagSet,
	gitHost *string,
	groupName *string,
	projectName *string,
	authKey *string,
	upgradePyFile *bool,
	versionService *version.Version) HandlerService {
	return &Handler{
		upgradeVersionCmd: upgradeVersionCmd,
		gitHost:           gitHost,
		groupName:         groupName,
		projectName:       projectName,
		authKey:           authKey,
		upgradePyFile:     upgradePyFile,
		versionService:    versionService,
	}
}
