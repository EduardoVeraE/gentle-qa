package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the braille ASCII art for the Gentle-QA owl mascot
// (inspired by the GitHub Copilot owl style).
var logoLines = []string{
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⢀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⢀⣾⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠙⠻⠿⠿⠿⠿⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠛⠿⠿⠿⠿⠟⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⡿⠿⠛⠉⠁⠀⠀⠀⠀⠉⠙⠻⢿⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⣴⣿⣿⠟⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠻⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⣼⣿⡿⠃⠀⠀⠀⢀⣴⣿⣷⡄⠀⠀⢠⣾⣿⣦⡀⠀⠀⠘⢿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢸⣿⣿⠁⠀⠀⠀⠀⣿⣿⣿⣿⡇⠀⠀⢸⣿⣿⣿⣿⠀⠀⠀⠈⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢸⣿⣿⠀⠀⠀⠀⠀⠙⠿⠿⠿⠁⠀⠀⠈⠿⠿⠿⠋⠀⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⢸⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⣶⣶⣶⡆⠀⠀⠀⠀⠀⠀⠀⣀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠈⢿⣿⣷⡀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⣠⣿⣿⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠈⠻⣿⣿⣦⣀⠀⠀⠀⠀⠘⠛⠛⠛⠋⠀⠀⠀⠀⣀⣴⣿⣿⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠿⣿⣿⣶⣤⣄⣀⣀⣀⣀⣀⣠⣤⣶⣿⣿⠿⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠙⠛⠿⠿⣿⣿⣿⠿⠿⠛⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
	"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
}

// gradientColors defines the top-to-bottom gradient for the Gentle-QA owl logo.
// Teal → Blue → Lavender → Mauve → Green — QE testing palette.
var gradientColors = []lipgloss.Color{
	ColorTeal,     // band 1 — ear tufts
	ColorBlue,     // band 2 — upper face
	ColorLavender, // band 3 — eyes
	ColorMauve,    // band 4 — lower face
	ColorGreen,    // band 5 — neck/body
}

// brandLine is rendered below the owl art: mascot text face + product name.
const brandLine = "  (^.^)  Gentle-QA"

// RenderLogo returns the Gentle-QA owl logo with gradient + brand line.
func RenderLogo() string {
	total := len(logoLines)
	if total == 0 {
		return ""
	}

	bands := len(gradientColors)
	var b strings.Builder

	for i, line := range logoLines {
		bandIdx := (i * bands) / total
		if bandIdx >= bands {
			bandIdx = bands - 1
		}
		style := lipgloss.NewStyle().Foreground(gradientColors[bandIdx])
		b.WriteString(style.Render(line))
		b.WriteByte('\n')
	}

	// Brand line: mascot face (Mauve) + product name (Lavender bold)
	mascotStyle := lipgloss.NewStyle().Foreground(ColorMauve).Bold(true)
	nameStyle := lipgloss.NewStyle().Foreground(ColorLavender).Bold(true)
	b.WriteString(mascotStyle.Render("  (^.^)"))
	b.WriteString(nameStyle.Render("  Gentle-QA"))

	return b.String()
}
