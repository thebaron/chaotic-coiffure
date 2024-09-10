package view

import (
	"github.com/charmbracelet/bubbles/spinner"
	sonic "github.com/delucks/go-subsonic"
)

type model struct {
	sonic    *sonic.Client
	spinner  spinner.Model
	quitting bool
	err      error
	msg      string
}
