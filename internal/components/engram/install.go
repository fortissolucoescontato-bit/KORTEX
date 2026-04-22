package engram

import (
	"github.com/fortissolucoescontato-bit/kortex/internal/installcmd"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentEngram)
}
