package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Padding(3, 3)

func (m model) View() string {
	return appStyle.Render(m.list.View())
}
