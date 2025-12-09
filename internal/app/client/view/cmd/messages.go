package cmd

import tea "github.com/charmbracelet/bubbletea"

type FtpCmd struct {
	Value string
}

func emitCmd(value string) tea.Cmd {
	return func () tea.Msg {
		return FtpCmd { Value: value }
	}
}
