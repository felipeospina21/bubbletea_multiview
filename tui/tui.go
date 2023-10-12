package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	useHighPerformanceRenderer              = false
	spinnerView                sessionState = iota
	listView
	childListView
	viewportView
)

type mainModel struct {
	state     sessionState
	spinner   spinner.Model
	list      list.Model
	childList list.Model
	viewport  viewportModel
	index     int
}

func NewModel(timeout time.Duration) mainModel {
	m := mainModel{state: listView}

	m.initSpinnerModel()
	m.initListModel()
	m.initChildListModel()
	m.initViewportModel()

	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.spinner.Tick)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.setListsViewsSize(msg, h, v)
		cmd := m.setViewportViewSize(msg, headerHeight, verticalMarginHeight)

		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.switchView()
		}

		switch m.state {
		// update whichever model is focused
		case spinnerView:
			if msg.String() == "n" {
				m.Next()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.Tick)
			}
		}
	}
	switch m.state {
	// update whichever model is focused
	case listView:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)

	case childListView:
		m.childList, cmd = m.childList.Update(msg)
		cmds = append(cmds, cmd)

	case viewportView:
		m.viewport.mod, cmd = m.viewport.mod.Update(msg)
		cmds = append(cmds, cmd)

	default:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	switch m.state {
	case listView:
		s += docStyle.Render(m.list.View())

	case childListView:
		s += docStyle.Render(m.childList.View())

	case viewportView:
		s += fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.mod.View(), m.footerView())

	default:
		s += docStyle.Render(m.spinner.View())
	}

	// model := m.currentFocusedModel()
	// s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m *mainModel) currentFocusedModel() string {
	switch m.state {
	case listView:
		return "list"

	case childListView:
		return "childListView"

	case viewportView:
		return "viewportView"

	default:
		return "spinner"
	}
}

func (m *mainModel) switchView() {
	switch m.state {

	case spinnerView:
		m.state = listView

	case listView:
		m.state = childListView

	case childListView:
		m.state = viewportView

	default:
		m.state = spinnerView
	}
}
