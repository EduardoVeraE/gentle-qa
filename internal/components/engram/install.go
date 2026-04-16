package engram

import (
	"github.com/EduardoVeraE/Gentle-QA/internal/installcmd"
	"github.com/EduardoVeraE/Gentle-QA/internal/model"
	"github.com/EduardoVeraE/Gentle-QA/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentEngram)
}
