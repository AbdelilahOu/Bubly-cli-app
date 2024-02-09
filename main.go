package main

import (
	"fmt"

	"github.com/AbdelilahOu/Bubly-cli-app/app"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(2)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
	Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	initialModel := app.AppModel{
		Choice:           0,
		Quitting:         false,
		History:          []string{"main"},
		Textarea:         ta,
		Viewport:         vp,
		Text:             "",
		IsTextAreaActive: false,
		IsUrlWritten:     false,
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
