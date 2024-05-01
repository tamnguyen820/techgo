package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tamnguyen820/techgo/internal/rss"
)

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type item struct {
	title, desc, publishedTime string
}

func (i item) Title() string         { return i.title }
func (i item) Description() string   { return i.desc }
func (i item) PublishedTime() string { return i.publishedTime }
func (i item) FilterValue() string   { return i.title }

func InitialModel() (model, error) {
	allFeeds, err := rss.FetchAllFeeds()
	if err != nil {
		return model{}, err
	}

	items := []list.Item{}
	for _, feed := range allFeeds {
		for _, entry := range feed.Feed.Items {
			items = append(items, item{
				title:         entry.Title,
				desc:          entry.Description,
				publishedTime: entry.PublishedParsed.String()})
		}
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "TechGo"

	return m, nil
}
