//go:build !js

package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textInput    textinput.Model
	path         string
	history      []string
	historyIndex int
	backupCmd    string
}

func New() Model {
	textInput := textinput.New()
	textInput.Focus()
	textInput.Placeholder = "Entrez une commande"
	textInput.CharLimit = 156
	textInput.Width = 40
	textInput.Prompt = "/ $ "

	return Model{textInput: textInput, history: []string{}, historyIndex: -1}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ChangeDirMsg:
		m.path = msg.To
		m.textInput.Prompt = m.path + "/ $ "
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:

			// get value and reset input
			value := m.textInput.Value()
			m.textInput.SetValue("")

			m.historyIndex = -1
			if len(m.history) == 0 || m.history[len(m.history)-1] != value {
				m.history = append(m.history, value)
			}

			// emit cmd
			return m, emitFtpCmd(value)

		case tea.KeyUp:

			if len(m.history) == 0 {
				break
			}

			if m.historyIndex == -1 {
				m.backupCmd = m.textInput.Value()
				m.historyIndex = len(m.history) - 1
			} else if m.historyIndex > 0 {
				m.historyIndex -= 1
			}

			m.textInput.SetValue(m.history[m.historyIndex])
			return m, cmd

		case tea.KeyDown:

			if len(m.history) == 0 || m.historyIndex == -1 {
				break
			}

			if m.historyIndex == len(m.history)-1 {
				m.historyIndex = -1
				m.textInput.SetValue(m.backupCmd)
				break
			}

			if len(m.history)-1 > m.historyIndex {
				m.historyIndex += 1
			}

			m.textInput.SetValue(m.history[m.historyIndex])
			return m, cmd

		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.textInput.View()
}
