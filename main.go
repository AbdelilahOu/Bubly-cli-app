package main

import (
	"fmt"

	"github.com/AbdelilahOu/Bubly-cli-app/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initialModel := app.AppModel{
		Choice:     0,
		Quitting:   false,
		History:    []string{"main"},
		ActiveView: "main",
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
