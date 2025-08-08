package app

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

// General stuff for styling the view
var (
	term      = termenv.EnvColorProfile()
	subtle    = makeFgStyle("241")
	dot       = colorFg(" • ", "236")
	help      = subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit") + dot + subtle("backspace: back")
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
)

var (
	TitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
	SuccessStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#16a34a", Dark: "#16a34a"}).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
	ErrorStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#b91c1c", Dark: "#b91c1c"}).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
)

type ViewsOptions struct {
	View        string
	ChoiceLabel string
}

type AppModel struct {
	Choice             int
	Quitting           bool
	History            []string
	Textarea           textarea.Model
	Text               string
	IsTextAreaActive   bool
	IsUrlWritten       bool
	PrintingIsDone     bool
	PrintingError      bool
	CancelBackgroudJob context.CancelFunc
	IsBackgroundJob    bool
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if m.IsTextAreaActive {
			if k == "esc" || k == "ctrl+c" {
				m.Quitting = true
				return m, tea.Quit
			}
		} else {
			if k == "q" || k == "esc" || k == "ctrl+c" {
				m.Quitting = true
				return m, tea.Quit
			}
		}
		if k == "backspace" && len(m.History) > 0 {
			if m.IsBackgroundJob {
				m.CancelBackgroudJob()
				m.IsBackgroundJob = false
			}
			if m.Textarea.Value() == "" {
				m.IsUrlWritten = false
				m.PrintingError = false
				m.PrintingIsDone = false
				m = removeFromHistory(m)
			} else {
				m.Text = m.Textarea.Value()[:len(m.Textarea.Value())-1]
			}
		}
	}

	return UpdateYoutube(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m AppModel) View() string {
	if m.Quitting {
		return "" + TitleStyle("See you later! 👋") + ""
	}
	s := YoutubeView(m)
	return indent.String(""+s+""+help, 2)
}

func destructureOptions(options []ViewsOptions, c int) []any {
	var choices []any
	for i, option := range options {
		choices = append(choices, checkbox(option.ChoiceLabel, c == i))
	}
	return choices
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Utils
func appendToHistory(m AppModel, s string) AppModel {
	m.History = append(m.History, s)
	return m
}
func removeFromHistory(m AppModel) AppModel {
	m.History = m.History[:len(m.History)-1]
	m.Choice = 0
	return m
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
