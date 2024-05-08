package tui

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Padding(2, 2)
var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render
var viewportStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

var glamourRenderer, _ = glamour.NewTermRenderer(
	glamour.WithAutoStyle(),
)

func (m model) View() string {
	switch m.viewMode {
	case FeedView:
		return m.showFeedView()
	case ArticleView:
		return m.showArticleView()
	}
	return ""
}

func (m model) showFeedView() string {
	return appStyle.Render(m.list.View())
}

func (m model) showArticleView() string {
	return m.articleViewPort.View()
}
