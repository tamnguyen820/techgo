package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func RefreshMsg() tea.Cmd {
	return func() tea.Msg {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
	}
}

type RefreshDone struct {
	articles []list.Item
}

func fetchAndSortArticles(m model) ([]list.Item, error) {
	allFeeds, err := m.rssService.FetchAllFeeds()
	if err != nil {
		return nil, err
	}

	var itemList customItemList
	for _, feed := range allFeeds {
		for _, entry := range feed.Feed.Items {
			var allAuthorNames []string
			for _, author := range entry.Authors {
				allAuthorNames = append(allAuthorNames, author.Name)
			}
			allAuthorsStr := strings.Join(allAuthorNames, ", ")

			itemList = append(itemList, customItem{
				title:         strings.TrimSpace(entry.Title),
				authors:       allAuthorsStr,
				feedName:      feed.FeedInfo.Name,
				url:           entry.Link,
				publishedTime: entry.PublishedParsed})
		}
	}
	// Sort the items by published time
	sort.Sort(itemList)

	tuiItems := make([]list.Item, len(itemList))
	for i, item := range itemList {
		tuiItems[i] = item
	}

	return tuiItems, nil
}

type customItemList []customItem

func (s customItemList) Len() int {
	return len(s)
}

func (s customItemList) Less(i, j int) bool {
	return s[i].publishedTime.After(*s[j].publishedTime)
}

func (s customItemList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
