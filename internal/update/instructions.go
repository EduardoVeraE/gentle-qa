package update

import (
	"github.com/EduardoVeraE/Gentle-QA/internal/system"
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
	case "gentle-qa":
		return gentleQAHint(profile)
	case "engram":
		return engramHint(profile)
	case "gga":
		return ggaHint(profile)
	case "opencode-subagent-statusline", "opencode-sdd-engram-manage":
		return "Restart/reload OpenCode; plugins are registered in ~/.config/opencode/tui.json"
	default:
		return ""
	}
}

func gentleQAHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade gentle-qa"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/EduardoVeraE/Gentle-QA/main/scripts/install.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/EduardoVeraE/Gentle-QA/main/scripts/install.ps1 | iex"
	default:
		return ""
	}
}

func engramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade engram"
	default:
		return "gentle-qa upgrade (downloads pre-built binary)"
	}
}

func ggaHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade gga"
	default:
		return "See https://github.com/Gentleman-Programming/gentleman-guardian-angel"
	}
}
