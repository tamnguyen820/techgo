package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tamnguyen820/techgo/internal/tui"
)

func main() {
	m, err := tui.InitialModel()
	if err != nil {
		fmt.Println("Error initializing model:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
