package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const useHighPerformanceRenderer = false
const (
	defaultTime              = time.Minute
	spinnerView sessionState = iota
	listView
	childListView
	viewportView
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	titleStyle   = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type viewportModel struct {
	mod     viewport.Model
	ready   bool
	content string
}
type mainModel struct {
	state     sessionState
	spinner   spinner.Model
	list      list.Model
	childList list.Model
	viewport  viewportModel
	index     int
}

func newModel(timeout time.Duration) mainModel {
	m := mainModel{state: listView}
	m.spinner = spinner.New()

	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "NNN", desc: "It's good on toast"},
		item{title: "AAA", desc: "It's good on toast"},
		item{title: "some", desc: "It's good on toast"},
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "list "

	childItems := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
	}

	m.childList = list.New(childItems, list.NewDefaultDelegate(), 0, 0)
	m.childList.Title = "child list"

	content, err := os.ReadFile("artichoke.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	m.viewport.content = string(content)
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
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.childList.SetSize(msg.Width-h, msg.Height-v)
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

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
	// model := m.currentFocusedModel()
	switch m.state {
	case listView:
		s += docStyle.Render(m.list.View())

	case childListView:
		s += docStyle.Render(m.childList.View())

	case viewportView:
		s += fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.mod.View(), m.footerView())
		// s += docStyle.Render(m.viewport.mod.View())

	default:
		s += docStyle.Render(m.spinner.View())
	}
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
