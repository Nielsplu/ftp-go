//go:build !js

// Rendu terminal (lipgloss) : réservé à la TUI native, le build WebAssembly
// affiche les logs via le pont JavaScript (cmd/wasm).
package log

import (
	"github.com/charmbracelet/lipgloss"
)

func renderWithColor(color, content string) string {
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return lineStyle.Render(content)
}

func (log LogOrigin) toStr() string {

	switch log {
	case Usr: return renderWithColor("#c2d2ffff", "[Usr]")
	case Sys: return renderWithColor("#f35656ff", "[Sys]")
	case Srv: return renderWithColor("#f6883eff", "[Srv]")
	}

	return "[Unknown]"
}
