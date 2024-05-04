package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tamnguyen820/techgo/internal/services"
)

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	rssService   *services.RSSService
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
		timePassedRounded = fmt.Sprintf("%dd ago", daysPassed)
	} else if hoursPassed > 0 {
		timePassedRounded = fmt.Sprintf("%dh ago", hoursPassed)
	} else {
		timePassedRounded = fmt.Sprintf("%dm ago", minutesPassed)
	}

	description := []string{}
	for _, s := range []string{i.feedName, i.authors, timePassedRounded} {
		if len(s) > 0 {
			description = append(description, s)
		}
	}

	return strings.Join(description, " | ")
}
func (i customItem) FilterValue() string { return i.title + i.feedName }

type listKeyMap struct {
	refresh key.Binding
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("TechGo"), RefreshMsg())
}

func NewModel(rssService *services.RSSService) (model, error) {
	var delegateKeys = newDelegateKeyMap()
	delegate := newItemDelegate(delegateKeys)
	var listKeys = newListKeyMap()

	tuiItems := []list.Item{}

	articleList := list.New(tuiItems, delegate, 0, 0)
	articleList.Title = "Tech News 📰"
	articleList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.refresh,
		}
	}
	articleList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.refresh,
		}
	}
	articleList.SetSpinner(spinner.Globe)

	m := model{
		list:         articleList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		rssService:   rssService,
	}

	return m, nil
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}
