package view

import tea "github.com/charmbracelet/bubbletea"

type NewProgressBarMsg struct {
	filename string
}

type RemoveProgressBarMsg struct {
	filename string
}

func removeProgressBar(filename string) tea.Msg {
	return RemoveProgressBarMsg{ filename }
}

func AddNewProgressBar(filename string) tea.Msg {
	return NewProgressBarMsg{ filename }
}

func listenForCmdIn(ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}