package app

import (
	"fmt"
	"os"
	"strings"
	"time"

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
		View:        "yt-download-subtitles",
		ChoiceLabel: "Download Youtube subtitles 📝",
	},
}

func UpdateYoutube(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	// Handle audio format selection if active
	if len(m.History) > 0 && m.History[0] == "yt-download-audio" && m.IsUrlWritten && m.AudioFormatSel != nil {
		// Create debug file
		debugFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if debugFile != nil {
			fmt.Fprintf(debugFile, "Handling audio format selection, IsUrlWritten: %t, AudioFormatSel != nil: %t\n", m.IsUrlWritten, m.AudioFormatSel != nil)
			if m.AudioFormatSel != nil {
				fmt.Fprintf(debugFile, "AudioFormatSel: Formats=%d, Choice=%d, Selected=%t\n", len(m.AudioFormatSel.Formats), m.AudioFormatSel.Choice, m.AudioFormatSel.Selected)
			}
			debugFile.Close()
		}
		
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j", "down":
				if len(m.AudioFormatSel.Formats) > m.AudioFormatSel.Choice+1 {
					m.AudioFormatSel.Choice++
					// Update page if necessary
					itemsPerPage := m.ItemsPerPage
					newPage := m.AudioFormatSel.Choice / itemsPerPage
					if newPage != m.Page {
						m.Page = newPage
					}
				}
				return m, nil
			case "k", "up":
				if m.AudioFormatSel.Choice > 0 {
					m.AudioFormatSel.Choice--
					// Update page if necessary
					itemsPerPage := m.ItemsPerPage
					newPage := m.AudioFormatSel.Choice / itemsPerPage
					if newPage != m.Page {
						m.Page = newPage
					}
				}
				return m, nil
			case "h", "left":
				// Navigate to previous page
				if m.Page > 0 {
					m.Page--
					// Update choice to first item on new page
					itemsPerPage := m.ItemsPerPage
					m.AudioFormatSel.Choice = m.Page * itemsPerPage
				}
				return m, nil
			case "l", "right":
				// Navigate to next page
				totalItems := len(m.AudioFormatSel.Formats)
				itemsPerPage := m.ItemsPerPage
				totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
				if m.Page < totalPages-1 {
					m.Page++
					// Update choice to first item on new page
					m.AudioFormatSel.Choice = m.Page * itemsPerPage
				}
				return m, nil
			case "enter":
				if !m.AudioFormatSel.Selected {
					m.AudioFormatSel.Selected = true
					m.AudioFormatSel.Downloading = true
					formatID := m.AudioFormatSel.Formats[m.AudioFormatSel.Choice].ID
					return m, m.downloadAudio(m.AudioFormatSel.URL, formatID)
				}
				return m, nil
			}
		case AudioDownloadMsg:
			m.AudioFormatSel.Downloading = false
			if msg.Error != "" {
				m.AudioFormatSel.Error = true
				m.AudioFormatSel.ErrMsg = msg.Error
			} else {
				m.AudioFormatSel.Done = true
			}
			return m, nil
		}
		return m, nil
	}

	// Handle video format selection if active
	if len(m.History) > 0 && m.History[0] == "yt-download-video" && m.IsUrlWritten && m.VideoFormatSel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j", "down":
				if len(m.VideoFormatSel.Formats) > m.VideoFormatSel.Choice+1 {
					m.VideoFormatSel.Choice++
					// Update page if necessary
					itemsPerPage := m.ItemsPerPage
					newPage := m.VideoFormatSel.Choice / itemsPerPage
					if newPage != m.Page {
						m.Page = newPage
					}
				}
				return m, nil
			case "k", "up":
				if m.VideoFormatSel.Choice > 0 {
					m.VideoFormatSel.Choice--
					// Update page if necessary
					itemsPerPage := m.ItemsPerPage
					newPage := m.VideoFormatSel.Choice / itemsPerPage
					if newPage != m.Page {
						m.Page = newPage
					}
				}
				return m, nil
			case "h", "left":
				// Navigate to previous page
				if m.Page > 0 {
					m.Page--
					// Update choice to first item on new page
					itemsPerPage := m.ItemsPerPage
					m.VideoFormatSel.Choice = m.Page * itemsPerPage
				}
				return m, nil
			case "l", "right":
				// Navigate to next page
				totalItems := len(m.VideoFormatSel.Formats)
				itemsPerPage := m.ItemsPerPage
				totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
				if m.Page < totalPages-1 {
					m.Page++
					// Update choice to first item on new page
					m.VideoFormatSel.Choice = m.Page * itemsPerPage
				}
				return m, nil
			case "enter":
				if !m.VideoFormatSel.Selected {
					m.VideoFormatSel.Selected = true
					m.VideoFormatSel.Downloading = true
					formatID := m.VideoFormatSel.Formats[m.VideoFormatSel.Choice].ID
					return m, m.downloadVideo(m.VideoFormatSel.URL, formatID)
				}
				return m, nil
			}
		case VideoDownloadMsg:
			m.VideoFormatSel.Downloading = false
			if msg.Error != "" {
				m.VideoFormatSel.Error = true
				m.VideoFormatSel.ErrMsg = msg.Error
			} else {
				m.VideoFormatSel.Done = true
			}
			return m, nil
		}
		return m, nil
	}

	// Handle subtitle language selection if active
	if len(m.History) > 0 && m.History[0] == "yt-download-subtitles" && m.IsUrlWritten && m.SubtitleSel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j", "down":
				if len(m.SubtitleSel.Languages) > m.SubtitleSel.Choice+1 {
					m.SubtitleSel.Choice++
					// Update page if necessary
					itemsPerPage := m.ItemsPerPage
					newPage := m.SubtitleSel.Choice / itemsPerPage
					if newPage != m.Page {
						m.Page = newPage
					}
				}
				return m, nil
			case "k", "up":
				if m.SubtitleSel.Choice > 0 {
					m.SubtitleSel.Choice--
					// Update page if necessary
					itemsPerPage := m.ItemsPerPage
					newPage := m.SubtitleSel.Choice / itemsPerPage
					if newPage != m.Page {
						m.Page = newPage
					}
				}
				return m, nil
			case "h", "left":
				// Navigate to previous page
				if m.Page > 0 {
					m.Page--
					// Update choice to first item on new page
					itemsPerPage := m.ItemsPerPage
					m.SubtitleSel.Choice = m.Page * itemsPerPage
				}
				return m, nil
			case "l", "right":
				// Navigate to next page
				totalItems := len(m.SubtitleSel.Languages)
				itemsPerPage := m.ItemsPerPage
				totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
				if m.Page < totalPages-1 {
					m.Page++
					// Update choice to first item on new page
					m.SubtitleSel.Choice = m.Page * itemsPerPage
				}
				return m, nil
			case "enter":
				if !m.SubtitleSel.Selected {
					m.SubtitleSel.Selected = true
					m.SubtitleSel.Downloading = true
					langCode := m.SubtitleSel.Languages[m.SubtitleSel.Choice].Code
					return m, m.downloadSubtitles(m.SubtitleSel.URL, langCode)
				}
				return m, nil
			}
		case SubtitleDownloadMsg:
			m.SubtitleSel.Downloading = false
			if msg.Error != "" {
				m.SubtitleSel.Error = true
				m.SubtitleSel.ErrMsg = msg.Error
			} else {
				m.SubtitleSel.Done = true
			}
			return m, nil
		}
		return m, nil
	}

	if len(m.History) > 0 {
		switch m.History[0] {
		case "yt-download-video":
			return UpdateDownloadVideo(msg, m)
		case "yt-download-audio":
			return UpdateDownloadAudio(msg, m)
		case "yt-download-subtitles":
			return UpdateDownloadSubtitles(msg, m)
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
	case AudioFormatMsg:
		// Create debug file
		debugFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if debugFile != nil {
			fmt.Fprintf(debugFile, "Main UpdateYoutube received AudioFormatMsg, Error: %s, Formats count: %d\n", msg.Error, len(msg.Formats))
			debugFile.Close()
		}
		
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
			m.PrintingError = true
		} else {
			m.AudioFormatSel = &AudioFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
			// Reset page when new formats are loaded
			m.Page = 0
		}
		return m, nil
	case VideoFormatMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
			m.PrintingError = true
		} else {
			m.VideoFormatSel = &VideoFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
			// Reset page when new formats are loaded
			m.Page = 0
		}
		return m, nil
	case SubtitleLangMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
			m.PrintingError = true
		} else {
			m.SubtitleSel = &SubtitleSelection{
				URL:       msg.URL,
				Languages: msg.Languages,
				Choice:    0,
			}
			// Reset page when new languages are loaded
			m.Page = 0
		}
		return m, nil
	case AudioDownloadMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.PrintingIsDone = true
		}
		return m, nil
	case VideoDownloadMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.PrintingIsDone = true
		}
		return m, nil
	case SubtitleDownloadMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.PrintingIsDone = true
		}
		return m, nil

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
		case "yt-download-subtitles":
			s.WriteString(DownloadSubtitlesView(m))
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
		// Show format selection if available
		if m.VideoFormatSel != nil {
			if m.VideoFormatSel.Error {
				s.WriteString(ErrorStyle("Error: " + m.VideoFormatSel.ErrMsg))
			} else if m.VideoFormatSel.Done {
				s.WriteString(SuccessStyle("Video downloaded successfully! Check assets folder"))
			} else if m.VideoFormatSel.Downloading {
				// Show downloading with spinner
				s.WriteString("📥 Downloading video")
				// Add a simple spinner animation
				spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
				frame := time.Now().UnixNano()/100000000 % int64(len(spinner))
				s.WriteString(" " + spinner[frame] + "\n")
				s.WriteString("This may take a few moments...")
			} else if len(m.VideoFormatSel.Formats) > 0 {
				s.WriteString("Select video format:\n\n")
				
				// Implement pagination
				totalItems := len(m.VideoFormatSel.Formats)
				itemsPerPage := m.ItemsPerPage
				totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
				currentPage := m.Page
				
				// Ensure page is within bounds
				if currentPage >= totalPages {
					currentPage = totalPages - 1
				}
				if currentPage < 0 {
					currentPage = 0
				}
				
				// Calculate start and end indices for current page
				startIdx := currentPage * itemsPerPage
				endIdx := startIdx + itemsPerPage
				if endIdx > totalItems {
					endIdx = totalItems
				}
				
				// Display items for current page
				for i := startIdx; i < endIdx; i++ {
					format := m.VideoFormatSel.Formats[i]
					cursor := "  "
					if m.VideoFormatSel.Choice == i {
						cursor = "> "
					}
					
					// Format the display line
					line := fmt.Sprintf("%s%s %s %s %s", 
						cursor, 
						videoQualityStyle(format.Quality), 
						videoFormatStyle(format.Format), 
						videoResolutionStyle(format.Resolution),
						videoFileSizeStyle(format.Filesize))
					
					s.WriteString(line + "\n")
				}
				
				// Display pagination info
				if totalPages > 1 {
					s.WriteString("\n")
					s.WriteString(fmt.Sprintf("Page %d of %d | ", currentPage+1, totalPages))
					if currentPage > 0 {
						s.WriteString("<-- Previous (h) ")
					}
					if currentPage < totalPages-1 {
						s.WriteString("Next (l) -->")
					}
				}
				
				s.WriteString("\n\n(Press ↑/↓ to select, Enter to download, h/l for pagination)")
			} else {
				s.WriteString("Loading available formats...")
			}
		} else {
			// Still showing the URL being processed
			s.WriteString("Fetching available video formats for: " + m.Text + "\n")
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
				// Fetch available video formats
				return m, m.fetchVideoFormats(m.Text)
			}
			return m, nil
		}
	case VideoFormatMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.VideoFormatSel = &VideoFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
		}
		return m, nil
	case VideoDownloadMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.PrintingIsDone = true
		}
		return m, nil
	}
	return m, tea.Batch(tiCmd)
}

// audio downloader view and update funcs
func DownloadAudioView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube audio \U0001F3B5"))
	s.WriteString("\n\n")

	// Create debug file
	debugFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if debugFile != nil {
		fmt.Fprintf(debugFile, "DownloadAudioView called, IsUrlWritten: %t\n", m.IsUrlWritten)
		if m.AudioFormatSel != nil {
			fmt.Fprintf(debugFile, "AudioFormatSel: Formats=%d, Choice=%d, Selected=%t, Downloading=%t, Done=%t, Error=%t\n", 
				len(m.AudioFormatSel.Formats), m.AudioFormatSel.Choice, m.AudioFormatSel.Selected, 
				m.AudioFormatSel.Downloading, m.AudioFormatSel.Done, m.AudioFormatSel.Error)
		} else {
			fmt.Fprintf(debugFile, "AudioFormatSel is nil\n")
		}
		debugFile.Close()
	}

	if m.IsUrlWritten {
		// Show format selection if available
		if m.AudioFormatSel != nil {
			if m.AudioFormatSel.Error {
				s.WriteString(ErrorStyle("Error: " + m.AudioFormatSel.ErrMsg))
			} else if m.AudioFormatSel.Done {
				s.WriteString(SuccessStyle("Audio downloaded successfully! Check assets folder"))
			} else if m.AudioFormatSel.Downloading {
				// Show downloading with spinner
				s.WriteString("🔊 Downloading audio")
				// Add a simple spinner animation
				spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
				frame := time.Now().UnixNano()/100000000 % int64(len(spinner))
				s.WriteString(" " + spinner[frame] + "\n")
				s.WriteString("This may take a few moments...")
			} else if len(m.AudioFormatSel.Formats) > 0 {
				s.WriteString("Select audio format:\n\n")
				
				// Implement pagination
				totalItems := len(m.AudioFormatSel.Formats)
				itemsPerPage := m.ItemsPerPage
				totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
				currentPage := m.Page
				
				// Ensure page is within bounds
				if currentPage >= totalPages {
					currentPage = totalPages - 1
				}
				if currentPage < 0 {
					currentPage = 0
				}
				
				// Calculate start and end indices for current page
				startIdx := currentPage * itemsPerPage
				endIdx := startIdx + itemsPerPage
				if endIdx > totalItems {
					endIdx = totalItems
				}
				
				// Display items for current page
				for i := startIdx; i < endIdx; i++ {
					format := m.AudioFormatSel.Formats[i]
					cursor := "  "
					if m.AudioFormatSel.Choice == i {
						cursor = "> "
					}
					
					// Format the display line
					line := fmt.Sprintf("%s%s %s %s", 
						cursor, 
						audioQualityStyle(format.Quality), 
						audioFormatStyle(format.Format), 
						audioFileSizeStyle(format.Filesize))
					
					s.WriteString(line + "\n")
				}
				
				// Display pagination info
				if totalPages > 1 {
					s.WriteString("\n")
					s.WriteString(fmt.Sprintf("Page %d of %d | ", currentPage+1, totalPages))
					if currentPage > 0 {
						s.WriteString("<-- Previous (h) ")
					}
					if currentPage < totalPages-1 {
						s.WriteString("Next (l) -->")
					}
				}
				
				s.WriteString("\n\n(Press \u2191/\u2193 to select, Enter to download, h/l for pagination)")
			} else {
				s.WriteString("Loading available formats...")
			}
		} else {
			// Still showing the URL being processed
			s.WriteString("Fetching available audio formats for: " + m.Text + "\n")
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
				// Fetch available audio formats
				// Create debug file
				debugFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if debugFile != nil {
					fmt.Fprintf(debugFile, "Calling fetchAudioFormats with URL: %s\n", m.Text)
					debugFile.Close()
				}
				return m, m.fetchAudioFormats(m.Text)
			}
			return m, nil
		}
	case AudioFormatMsg:
		// Create debug file
		debugFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if debugFile != nil {
			fmt.Fprintf(debugFile, "Received AudioFormatMsg, Error: %s, Formats count: %d\n", msg.Error, len(msg.Formats))
			debugFile.Close()
		}
		
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
			m.PrintingError = true
		} else {
			m.AudioFormatSel = &AudioFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
		}
		return m, nil
	}
	return m, tea.Batch(tiCmd)
}

// transcript downloader view and update funcs
func DownloadSubtitlesView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube subtitles \U0001F4DD"))
	s.WriteString("\n\n")

	if m.IsUrlWritten {
		// Show language selection if available
		if m.SubtitleSel != nil {
			if m.SubtitleSel.Error {
				s.WriteString(ErrorStyle("Error: " + m.SubtitleSel.ErrMsg))
			} else if m.SubtitleSel.Done {
				s.WriteString(SuccessStyle("Subtitles downloaded successfully! Check assets folder"))
			} else if m.SubtitleSel.Downloading {
				// Show downloading with spinner
				selectedLang := m.SubtitleSel.Languages[m.SubtitleSel.Choice].Name
				s.WriteString("📝 Downloading " + selectedLang + " subtitles")
				// Add a simple spinner animation
				spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
				frame := time.Now().UnixNano()/100000000 % int64(len(spinner))
				s.WriteString(" " + spinner[frame] + "\n")
				s.WriteString("This may take a few moments...")
			} else if len(m.SubtitleSel.Languages) > 0 {
				s.WriteString("Select subtitle language:\n\n")
				
				// Implement pagination
				totalItems := len(m.SubtitleSel.Languages)
				itemsPerPage := m.ItemsPerPage
				totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
				currentPage := m.Page
				
				// Ensure page is within bounds
				if currentPage >= totalPages {
					currentPage = totalPages - 1
				}
				if currentPage < 0 {
					currentPage = 0
				}
				
				// Calculate start and end indices for current page
				startIdx := currentPage * itemsPerPage
				endIdx := startIdx + itemsPerPage
				if endIdx > totalItems {
					endIdx = totalItems
				}
				
				// Display items for current page
				for i := startIdx; i < endIdx; i++ {
					lang := m.SubtitleSel.Languages[i]
					cursor := "  "
					if m.SubtitleSel.Choice == i {
						cursor = "> "
					}
					
					// Format the display line
					line := fmt.Sprintf("%s%s", cursor, subtitleLangStyle(lang.Name))
					
					s.WriteString(line + "\n")
				}
				
				// Display pagination info
				if totalPages > 1 {
					s.WriteString("\n")
					s.WriteString(fmt.Sprintf("Page %d of %d | ", currentPage+1, totalPages))
					if currentPage > 0 {
						s.WriteString("<-- Previous (h) ")
					}
					if currentPage < totalPages-1 {
						s.WriteString("Next (l) -->")
					}
				}
				
				s.WriteString("\n\n(Press ↑/↓ to select, Enter to download, h/l for pagination)")
			} else {
				s.WriteString("Loading available languages...")
			}
		} else {
			// Still showing the URL being processed
			s.WriteString("Fetching available subtitle languages for: " + m.Text + "\n")
		}
	} else {
		s.WriteString(m.Textarea.View())
	}
	return s.String()
}

func UpdateDownloadSubtitles(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
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
				// Fetch available subtitle languages
				return m, m.fetchSubtitleLanguages(m.Text)
			}
			return m, nil
		}
	case SubtitleLangMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.SubtitleSel = &SubtitleSelection{
				URL:       msg.URL,
				Languages: msg.Languages,
				Choice:    0,
			}
		}
		return m, nil
	case SubtitleDownloadMsg:
		if msg.Error != "" {
			m.PrintingError = true
			m.Warning = msg.Error
		} else {
			m.PrintingIsDone = true
		}
		return m, nil
	}
	return m, tea.Batch(tiCmd)
}
