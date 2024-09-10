package view

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/thebaron/chaotic-coiffure/pkg/client"
)

type model struct {
	spinner   spinner.Model
	quitting  bool
	err       error
	msg       string
	client    *client.Client
	clientsub chan client.Update
}
