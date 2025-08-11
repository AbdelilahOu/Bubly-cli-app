package app

import (
	"context"
	"fmt"

	"github.com/AbdelilahOu/Bubly-cli-app/types"
	"github.com/AbdelilahOu/Bubly-cli-app/utils"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

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
	WarningStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#f97316", Dark: "#f97316"}).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
)

var (
	audioQualityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	audioFormatStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 1).
				Render

	audioFileSizeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f97316")).
				Padding(0, 1).
				Render
)

var (
	videoQualityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	videoFormatStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 1).
				Render

	videoResolutionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f97316")).
				Padding(0, 1).
				Render

	videoFileSizeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2563eb")).
				Padding(0, 1).
				Render
)

var (
	subtitleLangStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	subtitleSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 1).
				Render
)

type ViewsOptions struct {
	View        string
	ChoiceLabel string
}

type AppModel struct {
	Choice               int
	Quitting             bool
	History              []string
	Textarea             textarea.Model
	Text                 string
	IsTextAreaActive     bool
	IsUrlWritten         bool
	PrintingIsDone       bool
	PrintingError        bool
	CancelBackgroudJob   context.CancelFunc
	IsBackgroundJob      bool
	Warning              string
	CheckingYtdlp        bool
	InstallingYtdlp      bool
	InstallationProgress int
	InstallationTotal    int
	InstallationMessage  string
	YtdlpInstalled       bool
	AudioFormatSel       *AudioFormatSelection
	VideoFormatSel       *VideoFormatSelection
	SubtitleSel          *SubtitleSelection
	Page                 int
	ItemsPerPage         int
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return types.CheckYtdlpMsg{Installed: utils.CheckYtdlp()}
		},
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	if m.CheckingYtdlp || m.InstallingYtdlp {
		return UpdateYtdlp(msg, m)
	}

	return UpdateYoutube(msg, m)
}

func (m AppModel) View() string {
	if m.Quitting {
		return "" + TitleStyle("See you later! 👋") + ""
	}

	if m.CheckingYtdlp || m.InstallingYtdlp {
		return YtdlpView(m)
	}

	s := YoutubeView(m)
	if m.Warning != "" {
		return indent.String(""+s+""+help+"\n"+WarningStyle(m.Warning), 2)
	}
	return indent.String(""+s+""+help, 2)
}

func UpdateYtdlp(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.CheckYtdlpMsg:
		m.YtdlpInstalled = msg.Installed
		if !m.YtdlpInstalled {
			m.CheckingYtdlp = true
		} else {
			m.CheckingYtdlp = false
		}
		return m, nil
	case types.YtdlpInstalledMsg:
		if msg.Err != nil {
			m.Warning = "Error installing yt-dlp: " + msg.Err.Error()
		} else {
			m.Warning = "yt-dlp installed successfully"
			m.YtdlpInstalled = true
		}
		m.CheckingYtdlp = false
		m.InstallingYtdlp = false
		return m, nil
	case types.ProgressMsg:
		m.InstallationProgress = msg.Progress
		m.InstallationTotal = msg.Total
		m.InstallationMessage = msg.Message
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.Choice == 0 {
				if m.CheckingYtdlp {
					m.InstallingYtdlp = true
					return m, utils.InstallYtdlp()
				}
			} else {
				if m.CheckingYtdlp {
					m.CheckingYtdlp = false
					m.Warning = "yt-dlp is not installed. Some features may not work."
				}
				return m, nil
			}
		case "up", "k":
			if m.Choice > 0 {
				m.Choice--
			}
		case "down", "j":
			if m.Choice < 1 {
				m.Choice++
			}
		}
	}
	return m, nil
}

func YtdlpView(m AppModel) string {
	var s string

	if m.CheckingYtdlp {
		if m.InstallingYtdlp {
			s += TitleStyle("Installing yt-dlp") + "\n\n"
			s += "Please wait while yt-dlp is being installed...\n"
		} else {
			s += TitleStyle("yt-dlp is not installed") + "\n\n"
			s += "Would you like to install it?\n\n"
		}
	} else {
		return ""
	}

	if !m.InstallingYtdlp {
		choices := []string{"Yes", "No"}
		for i, choice := range choices {
			s += checkbox(choice, m.Choice == i) + "\n"
		}
	}

	return indent.String(s+"\n"+help, 2)
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

func appendToHistory(m AppModel, s string) AppModel {
	m.History = append(m.History, s)
	return m
}
func removeFromHistory(m AppModel) AppModel {

	if len(m.History) > 0 {
		switch m.History[len(m.History)-1] {
		case "yt-download-audio":

			m.AudioFormatSel = nil
			m.IsUrlWritten = false
			m.Text = ""
			m.Textarea.Reset()
		case "yt-download-video":

			m.VideoFormatSel = nil
			m.IsUrlWritten = false
			m.Text = ""
			m.Textarea.Reset()
		case "yt-download-subtitles":

			m.SubtitleSel = nil
			m.IsUrlWritten = false
			m.Text = ""
			m.Textarea.Reset()
		}
	}

	m.History = m.History[:len(m.History)-1]
	m.Choice = 0

	m.PrintingError = false
	m.PrintingIsDone = false
	m.Warning = ""

	return m
}

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
