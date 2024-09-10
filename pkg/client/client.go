package client

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	sonic "github.com/delucks/go-subsonic"
	"github.com/thebaron/chaotic-coiffure/pkg/config"
)

type State int

type QuitMessage struct{}

const (
	ERRORED State = iota
	DISCONNECTED
	CONNECTED
	AUTHENTICATED
	FETCHED_PLAYLIST
	SHUTDOWN
)

type Update struct {
	State  State
	Msg    string
	Result interface{}
}

type Client struct {
	client *sonic.Client

	is_shutdown bool
	config      config.Config
	channel     chan Update
}

func NewClient(c config.Config, ch chan Update) *Client {
	return &Client{
		is_shutdown: false,
		config:      c,
		channel:     ch,
	}
}

func (c *Client) Shutdown() {
	c.is_shutdown = true
}

func (c *Client) Launch() tea.Cmd {

	return func() tea.Msg {
		for !c.is_shutdown {

			var baseUrl string = fmt.Sprintf("https://%s/", c.config.Server.Host)
			c.client = &sonic.Client{
				Client:     &http.Client{},
				BaseUrl:    baseUrl,
				User:       c.config.Server.User,
				ClientName: "chaotic-coiffure",
			}

			c.channel <- Update{State: CONNECTED, Msg: fmt.Sprintf("connected to %s", baseUrl)}

			err := c.client.Authenticate(c.config.Server.Password)

			if err != nil {
				c.channel <- Update{State: DISCONNECTED, Msg: fmt.Sprintf("authentication failed: %v", err)}
				continue
			}

			c.channel <- Update{State: AUTHENTICATED, Msg: fmt.Sprintf("conneceted to %s", baseUrl)}

			pls, err := c.client.GetPlaylists(nil)
			if err != nil {
				c.channel <- Update{State: ERRORED, Msg: fmt.Sprintf("failed to retrieve list of playlists: %v", err)}
				continue
			}

			c.channel <- Update{State: FETCHED_PLAYLIST, Msg: fmt.Sprintf("playlist retrieved, %d playlists", len(pls)), Result: pls}
			c.Shutdown()
		}
		return QuitMessage{}
	}
}
