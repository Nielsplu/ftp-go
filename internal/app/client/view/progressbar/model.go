package progressbar

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

type Model struct {
	Filename string
	progress progress.Model
	percent  float64
}

func New(filename string) Model {
	return Model{
		Filename: filename,
		progress: progress.New(progress.WithScaledGradient("#3e57f6ff", "#ee4545ff")),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		return m, tea.Quit

	case tea.WindowSizeMsg:

		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case ProgressMsg:

		m.percent = msg.Percent
		return m, nil

	}
	return m, nil
}

func (m Model) View() string {
	pad := strings.Repeat(" ", padding)
	return pad + m.Filename + "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n"
}
