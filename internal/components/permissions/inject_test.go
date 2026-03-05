package permissions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/agents/claude"
	"github.com/gentleman-programming/gentle-ai/internal/agents/cursor"
	"github.com/gentleman-programming/gentle-ai/internal/agents/gemini"
	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
	"github.com/gentleman-programming/gentle-ai/internal/agents/vscode"
)

func claudeAdapter() agents.Adapter   { return claude.NewAdapter() }
func opencodeAdapter() agents.Adapter { return opencode.NewAdapter() }
func geminiAdapter() agents.Adapter   { return gemini.NewAdapter() }
func cursorAdapter() agents.Adapter   { return cursor.NewAdapter() }
func vscodeAdapter() agents.Adapter   { return vscode.NewAdapter() }

func TestInjectOpenCodeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}

	path := filepath.Join(home, ".config", "opencode", "opencode.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file %q: %v", path, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"permission"`) {
		t.Fatal("opencode.json missing permission key")
	}
	if strings.Contains(text, `"permissions"`) {
		t.Fatal("opencode.json should use 'permission' (singular), not 'permissions'")
	}
	if !strings.Contains(text, `"bash"`) {
		t.Fatal("opencode.json permission missing bash section")
	}
	if !strings.Contains(text, `"read"`) {
		t.Fatal("opencode.json permission missing read section")
	}
}

func TestInjectDefaultModeIsValidForClaudeCode(t *testing.T) {
	home := t.TempDir()

	if _, err := Inject(home, claudeAdapter()); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	settingsPath := filepath.Join(home, ".claude", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read settings file %q: %v", settingsPath, err)
	}

	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("unmarshal settings json: %v", err)
	}

	permissionsNode, ok := settings["permissions"].(map[string]any)
	if !ok {
		t.Fatalf("permissions node missing or invalid: %#v", settings["permissions"])
	}

	mode, ok := permissionsNode["defaultMode"].(string)
	if !ok {
		t.Fatalf("defaultMode missing or not a string: %#v", permissionsNode["defaultMode"])
	}

	// Valid Claude Code permission modes per https://code.claude.com/docs/en/iam#permission-modes
	validModes := map[string]bool{
		"acceptEdits":       true,
		"bypassPermissions": true,
		"default":           true,
		"dontAsk":           true,
		"plan":              true,
	}

	if !validModes[mode] {
		t.Fatalf("defaultMode %q is not a valid Claude Code permission mode", mode)
	}
}

func TestInjectAddsEnvToDenyList(t *testing.T) {
	home := t.TempDir()

	if _, err := Inject(home, claudeAdapter()); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	settingsPath := filepath.Join(home, ".claude", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read settings file %q: %v", settingsPath, err)
	}

	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("unmarshal settings json: %v", err)
	}

	permissionsNode, ok := settings["permissions"].(map[string]any)
	if !ok {
		t.Fatalf("permissions node missing or invalid: %#v", settings["permissions"])
	}

	denyList, ok := permissionsNode["deny"].([]any)
	if !ok {
		t.Fatalf("deny list missing or invalid: %#v", permissionsNode["deny"])
	}

	for _, entry := range denyList {
		if value, ok := entry.(string); ok && value == ".env" {
			return
		}
	}

	t.Fatalf("deny list missing explicit .env rule: %#v", denyList)
}

func TestInjectGeminiCLIUsesApprovalMode(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, geminiAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() changed = false, expected true")
	}

	settingsPath := filepath.Join(home, ".gemini", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read settings file %q: %v", settingsPath, err)
	}

	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("unmarshal settings json: %v", err)
	}

	general, ok := settings["general"].(map[string]any)
	if !ok {
		t.Fatalf("general section missing or invalid: %#v", settings["general"])
	}

	mode, ok := general["defaultApprovalMode"].(string)
	if !ok {
		t.Fatalf("defaultApprovalMode missing or not a string: %#v", general["defaultApprovalMode"])
	}

	validModes := map[string]bool{
		"default":   true,
		"auto_edit": true,
		"plan":      true,
		"yolo":      true,
	}
	if !validModes[mode] {
		t.Fatalf("defaultApprovalMode %q is not a valid Gemini CLI mode", mode)
	}

	// Must not contain Claude Code-specific keys.
	if _, has := settings["permissions"]; has {
		t.Fatal("Gemini settings should not contain Claude Code 'permissions' key")
	}
}

func TestInjectVSCodeCopilotUsesAutoApprove(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, vscodeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() changed = false, expected true")
	}

	settingsPath := vscodeAdapter().SettingsPath(home)
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("read settings file %q: %v", settingsPath, err)
	}

	var settings map[string]any
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("unmarshal settings json: %v", err)
	}

	autoApprove, ok := settings["chat.tools.autoApprove"]
	if !ok {
		t.Fatal("chat.tools.autoApprove missing from VS Code settings")
	}
	if autoApprove != true {
		t.Fatalf("chat.tools.autoApprove = %v, want true", autoApprove)
	}

	// Must not contain Claude Code-specific keys.
	if _, has := settings["permissions"]; has {
		t.Fatal("VS Code settings should not contain Claude Code 'permissions' key")
	}
}

func TestInjectCursorIsNoOp(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, cursorAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if result.Changed {
		t.Fatal("Inject() changed = true for Cursor, expected no-op")
	}
}
