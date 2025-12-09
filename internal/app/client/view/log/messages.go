package log

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogOrigin int

func (log LogOrigin) toStr() string {
	lineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a3400eff"))
	switch log {
	case Usr: return lineStyle.Render("[Usr]")
	case Sys: return lineStyle.Render("[Sys]")
	case Srv: return lineStyle.Render("[Srv]")
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
