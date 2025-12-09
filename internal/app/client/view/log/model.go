package log

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Items  []LogMsg
	Width  int
	Height int
}

func New() Model {
	return Model{
		Items:  []LogMsg{},
	}
}

func (m *Model) AddItem(item LogMsg) {
	m.Items = append(m.Items, item)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LogMsg:

		m.AddItem(msg)

	case tea.WindowSizeMsg:

		m.Width = msg.Width
		m.Height = msg.Height - 10

	}

	return m, nil
}

func (m Model) View() string {
	
	content := ""
	
	for _, item := range m.Items {
		content += item.origin.toStr() + " : " + item.content + "\n"
	}

	return content
}