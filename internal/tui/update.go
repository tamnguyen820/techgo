package tui

import (
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
		case key.Matches(msg, m.keys.refresh):
			m.list.StartSpinner()
			m.list.NewStatusMessage(statusMessageStyle("Updating feed..."))

			return m, func() tea.Msg {
				newItems, err := fetchAndSortArticles(m)
				if err != nil {
					return m.list.NewStatusMessage(statusMessageStyle("Error fetching articles"))
				}
				return RefreshDone{articles: newItems}
			}
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
