package main

import (
	"log"
	"time"

	"example/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultTime = time.Minute

func main() {
	p := tea.NewProgram(tui.NewModel(defaultTime), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
