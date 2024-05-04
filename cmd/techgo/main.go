package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tamnguyen820/techgo/internal/services"
	"github.com/tamnguyen820/techgo/internal/tui"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "path to the config file")
	flag.Parse()

	rssService := services.NewRSSService(configFile)
	articleService := services.NewArticleService()

	m, err := tui.NewModel(rssService, articleService)
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
