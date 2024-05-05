package tui

import (
	"github.com/pkg/browser"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		// Don't match any key below if the list is filtering
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.openInBrowser):
			selectedItem := m.list.SelectedItem().(customItem)
			if err := browser.OpenURL(selectedItem.url); err != nil {
				return m, m.list.NewStatusMessage(statusMessageStyle("Unable to open link in browser"))
			}
			return m, m.list.NewStatusMessage(statusMessageStyle("Open link in browser"))
		case key.Matches(msg, m.keys.openInTerminal):
			selectedItem := m.list.SelectedItem().(customItem)
			if article, err := m.articleService.ExtractArticle(selectedItem.url); err != nil {
				return m, m.list.NewStatusMessage(statusMessageStyle("Error extracting article"))
			} else {
				return m, m.list.NewStatusMessage(statusMessageStyle(article.Title))
			}
		case key.Matches(msg, m.keys.refresh):
			return m, tea.Batch(
				m.list.StartSpinner(),
				m.list.NewStatusMessage(statusMessageStyle("Updating feed...")),
				func() tea.Msg {
					return StartRefresh{}
				},
			)
		}
	case StartRefresh:
		return m, func() tea.Msg {
			newItems, err := FetchAndSortArticles(m)
			if err != nil {
				return m.list.NewStatusMessage(statusMessageStyle("Error fetching articles"))
			}
			return RefreshDone{articles: newItems}
		}
	case RefreshDone:
		m.list.StopSpinner()
		statusCmd := m.list.NewStatusMessage(statusMessageStyle("Feed updated!"))
		setItemsCmd := m.list.SetItems(msg.articles)
		return m, tea.Batch(statusCmd, setItemsCmd)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
