package opencode

import (
	"os"
	"path/filepath"
	"testing"
)

const fixtureJSON = `{
  "anthropic": {
    "id": "anthropic",
    "env": ["ANTHROPIC_API_KEY"],
    "name": "Anthropic",
    "models": {
      "claude-sonnet-4-20250514": {
        "id": "claude-sonnet-4-20250514",
        "name": "Claude Sonnet 4",
        "family": "claude",
        "tool_call": true,
        "reasoning": false,
        "cost": {"input": 3.0, "output": 15.0},
        "limit": {"context": 200000, "output": 8192}
      },
      "claude-haiku-3-20240307": {
        "id": "claude-haiku-3-20240307",
        "name": "Claude Haiku 3",
        "family": "claude",
        "tool_call": true,
        "reasoning": false,
        "cost": {"input": 0.25, "output": 1.25},
        "limit": {"context": 200000, "output": 4096}
      }
    }
  },
  "openai": {
    "id": "openai",
    "env": ["OPENAI_API_KEY"],
    "name": "OpenAI",
    "models": {
      "gpt-4o": {
        "id": "gpt-4o",
        "name": "GPT-4o",
        "family": "gpt",
        "tool_call": true,
        "reasoning": false,
        "cost": {"input": 2.5, "output": 10.0},
        "limit": {"context": 128000, "output": 4096}
      },
      "o1-mini": {
        "id": "o1-mini",
        "name": "o1-mini",
        "family": "o1",
        "tool_call": false,
        "reasoning": true,
        "cost": {"input": 3.0, "output": 12.0},
        "limit": {"context": 128000, "output": 65536}
      }
    }
  },
  "noenv": {
    "id": "noenv",
    "env": [],
    "name": "No Env Provider",
    "models": {}
  }
}`

func writeFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "models.json")
	if err := os.WriteFile(path, []byte(fixtureJSON), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestLoadModels(t *testing.T) {
	path := writeFixture(t)

	providers, err := LoadModels(path)
	if err != nil {
		t.Fatalf("LoadModels() error = %v", err)
	}

	if len(providers) != 3 {
		t.Fatalf("provider count = %d, want 3", len(providers))
	}

	anthropic, ok := providers["anthropic"]
	if !ok {
		t.Fatal("missing anthropic provider")
	}
	if anthropic.Name != "Anthropic" {
		t.Fatalf("anthropic name = %q", anthropic.Name)
	}
	if len(anthropic.Models) != 2 {
		t.Fatalf("anthropic model count = %d, want 2", len(anthropic.Models))
	}
	if len(anthropic.Env) != 1 || anthropic.Env[0] != "ANTHROPIC_API_KEY" {
		t.Fatalf("anthropic env = %v", anthropic.Env)
	}
}

func TestLoadModelsFileNotFound(t *testing.T) {
	_, err := LoadModels("/nonexistent/models.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDetectAvailableProviders(t *testing.T) {
	path := writeFixture(t)
	providers, err := LoadModels(path)
	if err != nil {
		t.Fatalf("LoadModels() error = %v", err)
	}

	// Override envLookup for test.
	original := envLookup
	defer func() { envLookup = original }()

	envLookup = func(key string) string {
		if key == "ANTHROPIC_API_KEY" {
			return "sk-test"
		}
		return ""
	}

	available := DetectAvailableProviders(providers)
	if len(available) != 1 {
		t.Fatalf("available count = %d, want 1", len(available))
	}
	if available[0] != "anthropic" {
		t.Fatalf("available[0] = %q, want anthropic", available[0])
	}
}

func TestDetectAvailableProvidersSkipsNoEnv(t *testing.T) {
	path := writeFixture(t)
	providers, err := LoadModels(path)
	if err != nil {
		t.Fatalf("LoadModels() error = %v", err)
	}

	original := envLookup
	defer func() { envLookup = original }()

	// All env vars set.
	envLookup = func(key string) string {
		return "set"
	}

	available := DetectAvailableProviders(providers)
	// "noenv" provider has empty Env slice, should be excluded.
	for _, id := range available {
		if id == "noenv" {
			t.Fatal("provider with empty env should not be in available list")
		}
	}
}

func TestFilterModelsForSDD(t *testing.T) {
	path := writeFixture(t)
	providers, err := LoadModels(path)
	if err != nil {
		t.Fatalf("LoadModels() error = %v", err)
	}

	// OpenAI has 2 models, but o1-mini has tool_call=false.
	openai := providers["openai"]
	sddModels := FilterModelsForSDD(openai)
	if len(sddModels) != 1 {
		t.Fatalf("openai SDD model count = %d, want 1", len(sddModels))
	}
	if sddModels[0].ID != "gpt-4o" {
		t.Fatalf("filtered model = %q, want gpt-4o", sddModels[0].ID)
	}

	// Anthropic has 2 models, both with tool_call=true.
	anthropic := providers["anthropic"]
	sddModels = FilterModelsForSDD(anthropic)
	if len(sddModels) != 2 {
		t.Fatalf("anthropic SDD model count = %d, want 2", len(sddModels))
	}
}
