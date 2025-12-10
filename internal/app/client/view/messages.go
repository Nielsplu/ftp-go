package view

import tea "github.com/charmbracelet/bubbletea"

type NewProgressBarMsg struct {
	filename string
}

type RemoveProgressBarMsg struct {
	filename string
}

type outMsg struct {
	inner tea.Msg
} 

func removeProgressBar(filename string) tea.Msg {
	return RemoveProgressBarMsg{ filename }
}

func AddNewProgressBar(filename string) tea.Msg {
	return NewProgressBarMsg{ filename }
}

func listenForCmdIn(ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		// wrap cmd in outMsg, to be able to listen again
		inner := <-ch
		return outMsg { inner }
	}
}

type FatalError struct {
	message string
	err error
}

func NotifyFatalError(message string, err error) tea.Msg {
	return FatalError{ message, err }
}