package tui

import (
	"github.com/pkg/browser"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type delegateKeyMap struct {
	openInBrowser key.Binding
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var selectedItem *customItem
		if i, ok := m.SelectedItem().(customItem); ok {
			selectedItem = &i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.openInBrowser):
				if err := browser.OpenURL(selectedItem.url); err != nil {
					return m.NewStatusMessage(statusMessageStyle("Unable to open link in browser"))
				}
				return m.NewStatusMessage(statusMessageStyle("Open link in browser"))
			}
		}

		return nil
	}

	help := []key.Binding{keys.openInBrowser}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{d.openInBrowser}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.openInBrowser,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		openInBrowser: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open in browser"),
		),
	}
}
