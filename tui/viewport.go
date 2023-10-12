package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewportModel struct {
	mod     viewport.Model
	ready   bool
	content string
}

func (m *mainModel) initViewportModel() {
	content, err := os.ReadFile("artichoke.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	m.viewport.content = string(content)
}

func (m *mainModel) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.mod.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *mainModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.mod.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.mod.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *mainModel) setViewportViewSize(msg tea.WindowSizeMsg, headerHeight int, verticalMarginHeight int) tea.Cmd {
	if !m.viewport.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		m.viewport.mod = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.mod.YPosition = headerHeight
		m.viewport.mod.HighPerformanceRendering = useHighPerformanceRenderer
		m.viewport.mod.SetContent(m.viewport.content)
		m.viewport.ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.viewport.mod.YPosition = headerHeight + 1
	} else {
		m.viewport.mod.Width = msg.Width
		m.viewport.mod.Height = msg.Height - verticalMarginHeight
	}
	if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		// cmds = append(cmds, viewport.Sync(m.viewport.mod))
		return viewport.Sync(m.viewport.mod)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
