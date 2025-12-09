package progressbar

import tea "github.com/charmbracelet/bubbletea"

type ProgressMsg struct {
	Filename string
	Percent  float64
}

func SetProgressFor(
	filename string, 
	percent  float64,
) tea.Msg {
	return ProgressMsg{ filename, percent }	
}