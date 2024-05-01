package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Padding(2, 2)
var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render

func (m model) View() string {
	return appStyle.Render(m.list.View())
}
