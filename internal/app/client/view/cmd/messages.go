package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
)

type FtpCmdMsg struct {
	Value string
}

func emitFtpCmd(cmd string) tea.Cmd {
	return func () tea.Msg {
		return FtpCmdMsg { Value: cmd }
	}
}

type ChangeDirMsg struct {
    To string
}

func ChangeDir(to string) ChangeDirMsg {
    return ChangeDirMsg{To: to}
}