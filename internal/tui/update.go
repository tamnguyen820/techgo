package tui

import (
	"github.com/pkg/browser"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.articleViewPort.Width = msg.Width - h
		m.articleViewPort.Height = msg.Height - 5
		glamour, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.articleViewPort.Width),
		)
		if err == nil {
			glamourRenderer = glamour
			if m.viewMode == ArticleView {
				rendered, err := glamourRenderer.Render(m.currentArticle.Render())
				if err == nil {
					m.articleViewPort.SetContent(rendered)
				}
			}
		}
	case tea.KeyMsg:
		if msg.Type == tea.KeyEscape && m.viewMode == ArticleView {
			m.viewMode = FeedView
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
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
			if m.viewMode == ArticleView {
				return m, nil
			}
			selectedItem := m.list.SelectedItem().(customItem)
			if article, err := m.articleService.ExtractArticle(selectedItem.url); err != nil {
				return m, m.list.NewStatusMessage(statusMessageStyle("Error extracting article"))
			} else {
				m.viewMode = ArticleView
				m.currentArticle = &ArticleInfo{
					Title:         selectedItem.title,
					FeedName:      selectedItem.feedName,
					Authors:       selectedItem.authors,
					URL:           selectedItem.url,
					PublishedTime: selectedItem.publishedTime,
					CleanedText:   article.CleanedText,
				}
				rendered, err := glamourRenderer.Render(m.currentArticle.Render())
				if err != nil {
					return m, m.list.NewStatusMessage(statusMessageStyle("Error rendering article"))
				}
				m.articleViewPort.SetContent(rendered)
				m.articleViewPort.GotoTop()
				return m, nil
			}
		case key.Matches(msg, m.keys.refresh):
			if m.viewMode == ArticleView {
				return m, nil
			}
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
	switch m.viewMode {
	case FeedView:
		m.list, cmd = m.list.Update(msg)
	case ArticleView:
		m.articleViewPort, cmd = m.articleViewPort.Update(msg)
	}
	return m, cmd
}
