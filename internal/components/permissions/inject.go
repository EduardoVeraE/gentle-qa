package permissions

import (
	"fmt"
	"os"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

type InjectionResult struct {
	Changed bool
	Files   []string
}

// claudeCodeOverlayJSON uses Claude Code's permissions.defaultMode setting.
// Valid modes: "acceptEdits", "bypassPermissions", "default", "dontAsk", "plan".
var claudeCodeOverlayJSON = []byte("{\n  \"permissions\": {\n    \"defaultMode\": \"bypassPermissions\",\n    \"deny\": [\n      \"rm -rf /\",\n      \"sudo rm -rf /\",\n      \".env\"\n    ]\n  }\n}\n")

// openCodeOverlayJSON uses the OpenCode "permission" key with bash/read granularity.
var openCodeOverlayJSON = []byte("{\n  \"permission\": {\n    \"bash\": {\n      \"*\": \"allow\",\n      \"git commit *\": \"ask\",\n      \"git push *\": \"ask\",\n      \"git push\": \"ask\",\n      \"git push --force *\": \"ask\",\n      \"git rebase *\": \"ask\",\n      \"git reset --hard *\": \"ask\"\n    },\n    \"read\": {\n      \"*\": \"allow\",\n      \"*.env\": \"deny\",\n      \"*.env.*\": \"deny\",\n      \"**/.env\": \"deny\",\n      \"**/.env.*\": \"deny\",\n      \"**/secrets/**\": \"deny\",\n      \"**/credentials.json\": \"deny\"\n    }\n  }\n}\n")

// geminiCLIOverlayJSON uses Gemini CLI's general.defaultApprovalMode setting.
// Valid modes: "default", "auto_edit", "plan", "yolo".
var geminiCLIOverlayJSON = []byte("{\n  \"general\": {\n    \"defaultApprovalMode\": \"yolo\"\n  }\n}\n")

// vscodeCopilotOverlayJSON enables auto-approval for VS Code Copilot agent mode tools.
var vscodeCopilotOverlayJSON = []byte("{\n  \"chat.tools.autoApprove\": true\n}\n")

// overlayForAgent returns the correct permissions overlay for the given agent,
// or nil if the agent does not support settings-based permission configuration.
func overlayForAgent(agentID model.AgentID) []byte {
	switch agentID {
	case model.AgentClaudeCode:
		return claudeCodeOverlayJSON
	case model.AgentOpenCode:
		return openCodeOverlayJSON
	case model.AgentGeminiCLI:
		return geminiCLIOverlayJSON
	case model.AgentVSCodeCopilot:
		return vscodeCopilotOverlayJSON
	default:
		// Agents like Cursor manage permissions via UI/CLI, not settings.json.
		return nil
	}
}

func Inject(homeDir string, adapter agents.Adapter) (InjectionResult, error) {
	settingsPath := adapter.SettingsPath(homeDir)
	if settingsPath == "" {
		return InjectionResult{}, nil
	}

	overlay := overlayForAgent(adapter.Agent())
	if overlay == nil {
		return InjectionResult{}, nil
	}

	writeResult, err := mergeJSONFile(settingsPath, overlay)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{settingsPath}}, nil
}

func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	return filemerge.WriteFileAtomic(path, merged, 0o644)
}

var osReadFile = func(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read json file %q: %w", path, err)
	}

	return content, nil
}
