package screens

import (
	"strings"

	"github.com/EduardoVeraE/Gentle-QA/internal/model"
	"github.com/EduardoVeraE/Gentle-QA/internal/tui/styles"
)

// PresetOptions lists the presets shown in the TUI selector.
// Order is intentional: qe-sdet first (Gentle-QA default — full SDET stack),
// then the focal QE presets, then the inherited upstream presets, then custom.
func PresetOptions() []model.PresetID {
	return []model.PresetID{
		model.PresetQESDET,
		model.PresetQEFront,
		model.PresetQEPerf,
		model.PresetQEAPI,
		model.PresetFullGentleman,
		model.PresetEcosystemOnly,
		model.PresetMinimal,
		model.PresetCustom,
	}
}

var presetDescriptions = map[model.PresetID]string{
	model.PresetQESDET:        "Full SDET stack: all QE skills (Playwright + Selenium + k6 + Karate + a11y + ISTQB)",
	model.PresetQEFront:       "Frontend QE: Playwright stack (BDD/CLI/E2E/MCP/regression) + a11y + ISTQB",
	model.PresetQEPerf:        "Performance QE: k6 load testing + ISTQB foundation",
	model.PresetQEAPI:         "API QE: Karate DSL + api-testing (Playwright + REST Assured) + ISTQB",
	model.PresetFullGentleman: "Everything: memory, SDD, skills, docs, persona & security",
	model.PresetEcosystemOnly: "Core tools only: memory, SDD, skills & docs (no persona/security)",
	model.PresetMinimal:       "Just Engram persistent memory",
	model.PresetCustom:        "Choose components and skills manually; keep existing persona/settings unmanaged",
}

func RenderPreset(selected model.PresetID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Ecosystem Preset"))
	b.WriteString("\n\n")

	for idx, preset := range PresetOptions() {
		isSelected := preset == selected
		focused := idx == cursor
		b.WriteString(renderRadio(string(preset), isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+presetDescriptions[preset]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Back"}, cursor-len(PresetOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return b.String()
}
