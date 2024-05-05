package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
)

var (
	FastLineSpinner = spinner.Spinner{
		Frames: []string{"|", "/", "-", "\\"},
		FPS:    time.Second / 100,
	}
)
