package view

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	sonic "github.com/delucks/go-subsonic"
	"github.com/thebaron/chaotic-coiffure/pkg/config"
)

type errMsg error

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func InitialModel(c config.Config) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	client := sonic.Client{
		Client:     &http.Client{},
		BaseUrl:    fmt.Sprintf("https://%s/", c.Server.Host),
		User:       c.Server.User,
		ClientName: "chaotic-coiffure",
	}
	err := client.Authenticate(c.Server.Password)
	if err != nil {
		return model{err: errMsg(err)}
	}

	pls, err := client.GetPlaylists(nil)
	if err != nil {
		return model{err: errMsg(err)}
	}
	msg := pls[0].Name
	return model{spinner: s, sonic: &client, msg: msg}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
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
