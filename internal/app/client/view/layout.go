package view

import (
	"ftp/internal/app/client/view/cmd"
	"ftp/internal/app/client/view/log"
	"ftp/internal/app/client/view/progressbar"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	log          log.Model
	cmd          cmd.Model
	progressbars []progressbar.Model
	msgInChan 	 chan tea.Msg
	cmdOutChan   chan string
}

func New(msgInChan chan tea.Msg, cmdOutChan chan string) Model {
	return Model{
		log: log.New(),
		cmd: cmd.New(),
		progressbars: []progressbar.Model{},
		msgInChan: msgInChan,
		cmdOutChan: cmdOutChan,
	}
}

func (m Model) Init() tea.Cmd {
	return listenForCmdIn(m.msgInChan)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = []tea.Cmd{}
	var c tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		// update log
		m.log, c = m.log.Update(tea.WindowSizeMsg {
			Width: msg.Width,
			Height: msg.Height - (4 + len(m.progressbars) * 3),
		})
		cmds = append(cmds, c)

		// update cmd
		m.cmd, c = m.cmd.Update(msg)
		cmds = append(cmds, c)

		// update all progress bar
		for index, item := range(m.progressbars) { 
			m.progressbars[index], c = item.Update(msg)
			cmds = append(cmds, c)
		}

	case tea.KeyMsg:

	
		// check exit and clear 
		switch msg.String() {
        case "ctrl+c", "esc":
            return m, func() tea.Msg { return cmd.FtpCmd{ Value: "End"} }
		case "ctrl+l":
			m.log.Items = []log.LogMsg{}
			
        }

		// update cmd
		m.cmd, c = m.cmd.Update(msg)
		cmds = append(cmds, c)

	case FatalError:

		m.log, c = m.log.Update(log.AddLog(log.Sys, msg.message))
		cmds = append(cmds, c)
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(time.Second)
			println(msg.err.Error())
			return tea.Quit()
		})

	case log.LogMsg:

		// update log
		m.log, c = m.log.Update(msg)
		cmds = append(cmds, c)

	case NewProgressBarMsg:

		m.progressbars = append(m.progressbars, progressbar.New(
			msg.filename, 
		))

		m.log, c = m.log.Update(tea.WindowSizeMsg {
			Width: m.log.Width,
			Height: m.log.Height - 2,
		})
		cmds = append(cmds, c)

	case RemoveProgressBarMsg:

		// remove progress bar
		for index, item := range(m.progressbars) {
			if (item.Filename == msg.filename) {
				m.progressbars[index] = m.progressbars[len(m.progressbars) - 1]
				m.progressbars = m.progressbars[:len(m.progressbars)-1]

				m.log, c = m.log.Update(tea.WindowSizeMsg {
					Width: m.log.Width,
					Height: m.log.Height + 2,
				})
				cmds = append(cmds, c)
				break
			}
		}

	case progressbar.ProgressMsg:

		// update all progress bars
		for index, item := range(m.progressbars) {
			if (item.Filename == msg.Filename) {
				m.progressbars[index], c = item.Update(msg)
				cmds = append(cmds, c)	

				// remove progress bar 1s after
				if msg.Percent >= 1 {
					cmds = append(cmds, func() tea.Msg {
						time.Sleep(time.Second)
						return removeProgressBar(msg.Filename)
					})
				}
			}
		}

	case cmd.FtpCmd:

		if msg.Value == "clear" {
			m.log.Items = []log.LogMsg{}
			break
		}

		// send ftp commande 
		m.cmdOutChan <- msg.Value

	case outMsg:

		// execute acual cmd
		cmds = append(cmds, func() tea.Msg { return msg.inner })

		// start listening again
		cmds = append(cmds, listenForCmdIn(m.msgInChan))

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.log.View() + "\n"
	s += m.cmd.View()

	s += "\n---\n"
	for _, item := range m.progressbars {
		s += item.View()
	}

	return s
}