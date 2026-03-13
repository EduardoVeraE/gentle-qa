package opencode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// DefaultCachePath returns the default path to the OpenCode models cache file.
func DefaultCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cache", "opencode", "models.json")
}

// ModelCost holds the per-million-token pricing.
type ModelCost struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

// ModelLimit holds context and output token limits.
type ModelLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

// Model represents a single model within a provider.
type Model struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Family      string    `json:"family"`
	ToolCall    bool      `json:"tool_call"`
	Reasoning   bool      `json:"reasoning"`
	Cost        ModelCost `json:"cost"`
	Limit       ModelLimit `json:"limit"`
}

// Provider represents a model provider with its env vars and model catalog.
type Provider struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Env    []string         `json:"env"`
	Models map[string]Model `json:"models"`
}

// LoadModels parses the OpenCode models cache JSON file and returns providers keyed by ID.
func LoadModels(cachePath string) (map[string]Provider, error) {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("read models cache %q: %w", cachePath, err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse models cache: %w", err)
	}

	providers := make(map[string]Provider, len(raw))
	for id, providerJSON := range raw {
		var p Provider
		if err := json.Unmarshal(providerJSON, &p); err != nil {
			// Skip malformed providers.
			continue
		}
		p.ID = id
		providers[id] = p
	}

	return providers, nil
}

// envLookup is a package-level variable for testing.
var envLookup = os.Getenv

// DetectAvailableProviders returns provider IDs whose required env vars are all set.
// Results are sorted alphabetically.
func DetectAvailableProviders(providers map[string]Provider) []string {
	var available []string

	for id, provider := range providers {
		if len(provider.Env) == 0 {
			continue
		}
		allSet := true
		for _, envVar := range provider.Env {
			if envLookup(envVar) == "" {
				allSet = false
				break
			}
		}
		if allSet {
			available = append(available, id)
		}
	}

	sort.Strings(available)
	return available
}

// FilterModelsForSDD returns models from a provider that support tool_call (required for SDD phases).
// Results are sorted by model name.
func FilterModelsForSDD(provider Provider) []Model {
	var models []Model
	for _, m := range provider.Models {
		if m.ToolCall {
			models = append(models, m)
		}
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i].Name < models[j].Name
	})

	return models
}

// SDDPhases returns the ordered list of SDD phase sub-agent names.
func SDDPhases() []string {
	return []string{
		"sdd-init",
		"sdd-explore",
		"sdd-propose",
		"sdd-spec",
		"sdd-design",
		"sdd-tasks",
		"sdd-apply",
		"sdd-verify",
		"sdd-archive",
	}
}
