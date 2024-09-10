package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thebaron/chaotic-coiffure/pkg/client"
	"github.com/thebaron/chaotic-coiffure/pkg/config"
)

type errMsg error

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

// A command that waits for the activity on a channel.
func waitForActivity(m model) tea.Cmd {
	return func() tea.Msg {
		var ret client.Update = <-m.clientsub
		m.msg = ret.Msg
		return ret
	}
}

func InitialModel(c config.Config) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := model{spinner: s, msg: "starting up", clientsub: make(chan client.Update)}
	m.client = client.NewClient(c, m.clientsub)

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForActivity(m), // wait for activity
		m.client.Launch(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit

		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	case client.QuitMessage:
		m.quitting = true
		return m, nil

	case client.Update:
		m.spinner, _ = m.spinner.Update(msg)
		return m, waitForActivity(m) // wait for next event

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s %s %s\n\n", m.spinner.View(), m.msg, quitKeys.Help().Desc)
	if m.quitting {
		return str + "\n"
	}
	return str
}
