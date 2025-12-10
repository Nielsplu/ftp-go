package log

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogOrigin int

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

const (
	Usr LogOrigin = iota
	Sys
	Srv
)

type LogMsg struct {
	origin  LogOrigin
	content string
}

func AddLog(origin  LogOrigin, content string) tea.Msg {
	return LogMsg{ origin, content }
}
