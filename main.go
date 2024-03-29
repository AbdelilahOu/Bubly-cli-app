package main

import (
	"fmt"

	"github.com/AbdelilahOu/Bubly-cli-app/app"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	ta := textarea.New()
	ta.Placeholder = "Pass in a url..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(2)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	initialModel := app.AppModel{
		Choice:           0,
		Quitting:         false,
		History:          []string{"main"},
		Textarea:         ta,
		Text:             "",
		IsTextAreaActive: false,
		IsUrlWritten:     false,
		PrintingIsDone:   false,
		PrintingError:    false,
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
