package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh"

func main() {
	p := tea.NewProgram(initApplicationStateModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
