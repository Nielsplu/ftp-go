package cmd

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textInput textinput.Model
}

func New() Model {
	textInput := textinput.New()
	textInput.Focus()
	textInput.Placeholder = "Entrez une commande"
	textInput.CharLimit = 156
	textInput.Width = 40

	return Model{ textInput }
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if (msg.Type == tea.KeyEnter) {

			// get value and reset input
			value := m.textInput.Value()
			m.textInput.SetValue("")
			
			// emit cmd
			return m, emitCmd(value)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.textInput.View()
}