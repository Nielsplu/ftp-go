package log

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LogOrigin int

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
