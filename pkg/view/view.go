package view

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	sonic "github.com/delucks/go-subsonic"

	// "github.com/thebaron/chaotic-coiffure/pkg/client"
	"github.com/thebaron/chaotic-coiffure/pkg/config"
)

type errMsg error
type State int

var message string = "default"

const (
	ERRORED State = iota
	DISCONNECTED
	CONNECTED
	AUTHENTICATED
	FETCHED_PLAYLIST
	SHUTDOWN
)

type updateMsg struct {
	state  State
	msg    string
	result interface{}
}

type model struct {
	spinner  spinner.Model
	quitting bool
	err      error
	msg      string
	c        config.Config
	status   updateMsg
	client   sonic.Client
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func InitialModel(c config.Config) *model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &model{spinner: s, msg: "disconnected", c: c}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.Connect(fmt.Sprintf("https://%s/", m.c.Server.Host)),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case updateMsg:
		m.spinner, _ = m.spinner.Update(msg)
		if msg.state == ERRORED {
			m.err = fmt.Errorf("%s", msg.msg)
			return m, nil
		}

		if msg.state == CONNECTED {
			m.msg = "Connected!"
			return m, m.GetPlaylist()
		} else if msg.state == FETCHED_PLAYLIST {
			var pls []*sonic.Playlist = msg.result.([]*sonic.Playlist)
			m.msg = fmt.Sprintf("playlist retrieved: %d total, 0=%s ", len(pls), pls[0].Name)
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	t := time.Now()
	str := fmt.Sprintf("\n\n   %s %d %s %s\n\n", m.spinner.View(), t.Second(), m.msg, quitKeys.Help().Desc)
	if m.quitting {
		return str + "\n"
	}
	return str
}

func (m *model) Connect(baseUrl string) tea.Cmd {
	return func() tea.Msg {
		m.client = sonic.Client{
			Client:       &http.Client{},
			BaseUrl:      baseUrl,
			User:         m.c.Server.User,
			PasswordAuth: true,
			ClientName:   "chaotic-coiffure",
		}

		err := m.client.Authenticate(m.c.Server.Password)
		if err != nil {
			return updateMsg{state: ERRORED, msg: fmt.Sprintf("authentication failed: %v", err)}
		}

		return updateMsg{state: CONNECTED, msg: fmt.Sprintf("connected to %s", baseUrl)}
	}
}

func (m *model) GetPlaylist() tea.Cmd {

	return func() tea.Msg {

		time.Sleep(1 * time.Second)
		pls, err := m.client.GetPlaylists(nil)
		if err != nil {
			return updateMsg{state: ERRORED, msg: fmt.Sprintf("failed to retrieve list of playlists: %v", err)}
		}

		return updateMsg{state: FETCHED_PLAYLIST, msg: fmt.Sprintf("playlist retrieved, %d playlists", len(pls)), result: pls}
	}
}
