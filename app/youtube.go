package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var YoutubeOptions = []ViewsOptions{
	{
		View:        "yt-download-video",
		ChoiceLabel: "Download Youtube video 📥",
	},
	{
		View:        "yt-download-audio",
		ChoiceLabel: "Download Youtube audio 🎵",
	},
	{
		View:        "yt-download-transcript",
		ChoiceLabel: "Download Youtube transcript 📝",
	},
}

func UpdateYoutube(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	if len(m.History) > 0 {
		switch m.History[0] {
		case "yt-download-video":
			return UpdateDownloadVideo(msg, m)
		case "yt-download-audio":
			return UpdateDownloadAudio(msg, m)
		case "yt-download-transcript":
			return UpdateDownloadTranscript(msg, m)
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(YoutubeOptions) > m.Choice+1 {
				m.Choice++
			}
		case "k", "up":
			if m.Choice > 0 {
				m.Choice--
			}
		case "enter":
			m.IsTextAreaActive = true
			m = appendToHistory(m, YoutubeOptions[m.Choice].View)
			return m, nil
		}

	}
	return m, nil
}

func YoutubeView(m AppModel) string {
	c := m.Choice
	var s strings.Builder

	if len(m.History) > 0 {
		switch m.History[0] {
		case "yt-download-video":
			s.WriteString(DownloadVideoView(m))
		case "yt-download-audio":
			s.WriteString(DownloadAudioView(m))
		case "yt-download-transcript":
			s.WriteString(DownloadTranscriptView(m))
		}
		s.WriteString("\n\n")
	} else {
		s.WriteString(TitleStyle("What youtube tools do you wanna use? 🔨"))
		s.WriteString("\n\n")

		choices := fmt.Sprintf(
			strings.Repeat("%s\n", len(YoutubeOptions)),
			destructureOptions(YoutubeOptions, c)...,
		)
		s.WriteString(choices)
		s.WriteString("\n")
	}

	return s.String()
}

// video downloader view and update funcs
func DownloadVideoView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube video 📥"))
	s.WriteString("\n\n")

	if m.IsUrlWritten {
		if m.PrintingError {
			s.WriteString(ErrorStyle("An error accured while downloading video"))
		} else if m.PrintingIsDone {
			s.WriteString(SuccessStyle("Downloading video done check assets folder"))
		} else {
			s.WriteString("Downloading video from : " + m.Text)
		}
	} else {
		s.WriteString(m.Textarea.View())
	}
	return s.String()
}

func UpdateDownloadVideo(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)
	m.Textarea, tiCmd = m.Textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.IsUrlWritten {
				m.Text = m.Textarea.Value()
				m.Textarea.Reset()
				m.IsUrlWritten = true
				m.IsTextAreaActive = false
				// a function to download video should be called here
			}
			return m, nil
		}
	}
	return m, tea.Batch(tiCmd)
}

// audio downloader view and update funcs
func DownloadAudioView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube audio 🎵"))
	s.WriteString("\n\n")

	if m.IsUrlWritten {
		if m.PrintingError {
			s.WriteString(ErrorStyle("An error accured while downloading audio"))
		} else if m.PrintingIsDone {
			s.WriteString(SuccessStyle("Downloading audio done check assets folder"))
		} else {
			s.WriteString("Downloading audio from : " + m.Text)
		}
	} else {
		s.WriteString(m.Textarea.View())
	}
	return s.String()
}

func UpdateDownloadAudio(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)
	m.Textarea, tiCmd = m.Textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.IsUrlWritten {
				m.Text = m.Textarea.Value()
				m.Textarea.Reset()
				m.IsUrlWritten = true
				m.IsTextAreaActive = false
				// a function to download audio should be called here
			}
			return m, nil
		}
	}
	return m, tea.Batch(tiCmd)
}

// transcript downloader view and update funcs
func DownloadTranscriptView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube transcript 📝"))
	s.WriteString("\n\n")

	if m.IsUrlWritten {
		if m.PrintingError {
			s.WriteString(ErrorStyle("An error accured while downloading transcript"))
		} else if m.PrintingIsDone {
			s.WriteString(SuccessStyle("Downloading transcript done check assets folder"))
		} else {
			s.WriteString("Downloading transcript from : " + m.Text)
		}
	} else {
		s.WriteString(m.Textarea.View())
	}
	return s.String()
}

func UpdateDownloadTranscript(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)
	m.Textarea, tiCmd = m.Textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.IsUrlWritten {
				m.Text = m.Textarea.Value()
				m.Textarea.Reset()
				m.IsUrlWritten = true
				m.IsTextAreaActive = false
				// a function to download transcript should be called here
			}
			return m, nil
		}
	}
	return m, tea.Batch(tiCmd)
}
