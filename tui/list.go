package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m *mainModel) initListModel() {
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "NNN", desc: "It's good on toast"},
		item{title: "AAA", desc: "It's good on toast"},
		item{title: "some", desc: "It's good on toast"},
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "list "
}

func (m *mainModel) initChildListModel() {
	childItems := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
	}

	m.childList = list.New(childItems, list.NewDefaultDelegate(), 0, 0)
	m.childList.Title = "child list"
}

func (m *mainModel) setListsViewsSize(msg tea.WindowSizeMsg, h int, v int) {
	m.list.SetSize(msg.Width-h, msg.Height-v)
	m.childList.SetSize(msg.Width-h, msg.Height-v)
}
