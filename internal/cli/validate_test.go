package cli

import (
	"testing"

	"github.com/EduardoVeraE/Gentle-QA/internal/model"
)

// TestNormalizePresetDefaultIsQESDET is a fork-defense test for Gentle-QA.
//
// Gentle-QA ships qe-sdet as the default preset (the QA-focused fork bundles
// the full SDET stack out of the box). Upstream gentle-ai uses a different
// default — every upstream-sync MUST preserve this fork-specific override.
//
// If this test fails after a `scripts/sync-upstream.sh` merge, an upstream
// commit overwrote the default in `internal/cli/validate.go`. Re-apply the
// rewrite there before completing the merge. See
// `internal/assets/skills/upstream-sync/SKILL.md` § Fork-specific defaults.
func TestNormalizePresetDefaultIsQESDET(t *testing.T) {
	got, err := normalizePreset("")
	if err != nil {
		t.Fatalf("normalizePreset(\"\") error = %v, want nil", err)
	}
	if got != model.PresetQESDET {
		t.Fatalf("normalizePreset(\"\") = %q, want %q (fork default)", got, model.PresetQESDET)
	}
}

// TestNormalizePresetAcceptsAllKnownPresets ensures every PresetID constant
// is accepted by normalizePreset. If a new preset is added in types.go but
// not wired through validate.go, this test fails.
func TestNormalizePresetAcceptsAllKnownPresets(t *testing.T) {
	cases := []model.PresetID{
		model.PresetFullGentleman,
		model.PresetEcosystemOnly,
		model.PresetMinimal,
		model.PresetCustom,
		model.PresetQEFront,
		model.PresetQEPerf,
		model.PresetQEAPI,
		model.PresetQESDET,
	}
	for _, want := range cases {
		got, err := normalizePreset(string(want))
		if err != nil {
			t.Errorf("normalizePreset(%q) error = %v, want nil", want, err)
			continue
		}
		if got != want {
			t.Errorf("normalizePreset(%q) = %q, want %q", want, got, want)
		}
	}
}

func TestNormalizePresetRejectsUnknown(t *testing.T) {
	if _, err := normalizePreset("not-a-real-preset"); err == nil {
		t.Fatal("normalizePreset(\"not-a-real-preset\") error = nil, want error")
	}
}
