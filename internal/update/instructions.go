package update

import (
	"github.com/fortissolucoescontato-bit/kortex/internal/system"
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
	case "kortex":
		return kortexHint(profile)
	case "kortex-engram":
		return KortexEngramHint(profile)
	default:
		return ""
	}
}

func kortexHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade kortex"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/fortissolucoescontato-bit/kortex/main/scripts/install.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/fortissolucoescontato-bit/kortex/main/scripts/install.ps1 | iex"
	default:
		return "Veja https://github.com/fortissolucoescontato-bit/kortex"
	}
}

func KortexEngramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade KortexEngram"
	default:
		return "kortex upgrade (baixa o binário pré-compilado)"
	}
}
