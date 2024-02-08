package app

import (
	"fmt"

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
)

type ViewsOptions struct {
	View        string
	ChoiceLabel string
}

type AppModel struct {
	Choice     int
	Quitting   bool
	History    []string
	ActiveView string
}

var MainOptions = []ViewsOptions{
	{
		View:        "youtube",
		ChoiceLabel: "Youtube tools 📺",
	},
	{
		View:        "scraping",
		ChoiceLabel: "Web scraping tools 🕸️",
	},
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		if k == "backspace" && len(m.History) > 1 {
			m = removeFromHistory(m)
		}
	}

	switch m.ActiveView {
	case "youtube":
		return UpdateYoutube(msg, m)
	}
	return UpdateMain(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m AppModel) View() string {
	var s string
	if m.Quitting {
		return "\n  " + TitleStyle("See you later!") + "\n\n"

	}
	switch m.ActiveView {
	case "youtube":
		s = YoutubeView(m)
	default:
		s = MainView(m)
	}

	return indent.String("\n"+s+"\n\n"+help, 2)
}

// Update loop for the first view where you're choosing a task.
func UpdateMain(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(MainOptions) > m.Choice+1 {
				m.Choice++
			}
		case "k", "up":
			if m.Choice > 0 {
				m.Choice--
			}
		case "enter":
			m = appendToHistory(m, MainOptions[m.Choice].View)
			return m, nil
		}

	}
	return m, nil
}

// The first view, where you're choosing a task
func MainView(m AppModel) string {
	c := m.Choice
	tpl := TitleStyle("What tools do you wanna use? 🔨") + "\n\n%s"
	choices := fmt.Sprintf(
		"%s\n%s\n",
		checkbox(MainOptions[0].ChoiceLabel, c == 0),
		checkbox(MainOptions[1].ChoiceLabel, c == 1),
	)

	return fmt.Sprintf(tpl, choices)
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
	m.ActiveView = s
	return m
}
func removeFromHistory(m AppModel) AppModel {
	m.History = m.History[:len(m.History)-1]
	m.ActiveView = m.History[len(m.History)-1]
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
