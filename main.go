package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	tea "github.com/charmbracelet/bubbletea"

	"carpet-cleaner/tui"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file:", err)
	}

	owner := os.Getenv("GITHUB_OWNER")
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		log.Fatal("Error: GITHUB_TOKEN is required.")
	}
	if owner == "" {
		log.Fatal("Error: GITHUB_OWNER is required in .env for deletion to work.")
	}

	model := tui.NewModel(owner, token)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}