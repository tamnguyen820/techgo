package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tamnguyen820/techgo/internal/rss"
)

type model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

type customItem struct {
	title         string
	authors       string
	feedName      string
	url           string
	publishedTime *time.Time
}

func (i customItem) Title() string { return i.title }
func (i customItem) Description() string {
	// Calculate how much time has passed since the published time
	timePassed := time.Since(*i.publishedTime)
	// Round to the nearest minute, hour, or day
	minutesPassed := int(timePassed.Minutes())
	hoursPassed := int(timePassed.Hours())
	daysPassed := int(timePassed.Hours() / 24)
	var timePassedRounded string
	if daysPassed > 0 {
		timePassedRounded = fmt.Sprintf("%dd", daysPassed)
	} else if hoursPassed > 0 {
		timePassedRounded = fmt.Sprintf("%dh", hoursPassed)
	} else {
		timePassedRounded = fmt.Sprintf("%dm", minutesPassed)
	}

	return fmt.Sprintf("%s | %s | %s ago", i.feedName, i.authors, timePassedRounded)
}
func (i customItem) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func NewModel() (model, error) {
	var delegateKeys = newDelegateKeyMap()
	delegate := newItemDelegate(delegateKeys)

	allFeeds, err := rss.FetchAllFeeds()
	if err != nil {
		return model{}, err
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
				title:         entry.Title,
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

	m := model{
		list:         list.New(tuiItems, delegate, 0, 0),
		delegateKeys: delegateKeys,
	}
	m.list.Title = "TechGo"

	return m, nil
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
