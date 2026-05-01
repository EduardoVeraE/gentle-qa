package update

import (
	"fmt"
	"strings"

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
		return "gentle-qa upgrade updates ~/.config/opencode npm deps, clears this plugin's @latest cache, then requires OpenCode restart/reload"
	default:
		return ""
	}
}

func openCodeRegisteredNotMaterializedHint(tool ToolInfo) string {
	pkg := strings.TrimSpace(tool.NpmPackage)
	if pkg == "" {
		pkg = tool.Name
	}
	return fmt.Sprintf("registered in ~/.config/opencode/tui.json; pending npm dependency materialization for %s. Run gentle-qa upgrade to install/update ~/.config/opencode dependencies, then restart or reload OpenCode; if it stays pending, check OpenCode logs for package or peer dependency errors.", pkg)
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
