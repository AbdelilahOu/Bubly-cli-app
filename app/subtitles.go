package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SubtitleLanguage struct {
	Code string
	Name string
}

type SubtitleSelection struct {
	URL         string
	Languages   []SubtitleLanguage
	Choice      int
	Selected    bool
	Downloading bool
	Done        bool
	Error       bool
	ErrMsg      string
}

func (m AppModel) fetchSubtitleLanguages(url string) tea.Cmd {
	return func() tea.Msg {
		// Create assets directory if it doesn't exist
		os.MkdirAll("assets", 0755)

		// Run yt-dlp to get subtitle information
		var path, ffmpegPath string
		if isWindows() {
			path = "bin/yt-dlp.exe"
			ffmpegPath = "bin/ffmpeg.exe"
		} else {
			path = "bin/yt-dlp"
			ffmpegPath = "bin/ffmpeg"
		}

		// Create or open log file
		logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return SubtitleLangMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		// Create buffers to capture output
		var outBuf, errBuf strings.Builder
		
		// Check if ffmpeg exists
		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil

		// Prepare arguments
		var args []string
		args = append(args, "--list-subs", url)
		
		// Add ffmpeg location if it exists
		if useFfmpeg {
			args = append(args, "--ffmpeg-location", ffmpegPath)
		}
		
		cmd := exec.Command(path, args...)
		cmd.Stdout = io.MultiWriter(&outBuf, logFile)
		cmd.Stderr = io.MultiWriter(&errBuf, logFile)
		
		err = cmd.Run()
		
		// Write debug info
		debugFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if debugFile != nil {
			defer debugFile.Close()
			fmt.Fprintf(debugFile, "Fetching subtitle languages for URL: %s\n", url)
			fmt.Fprintf(debugFile, "Command executed, err: %v\n", err)
			if err == nil {
				output := outBuf.String()
				fmt.Fprintf(debugFile, "Output length: %d\n", len(output))
				// Only write first 1000 chars to debug to avoid huge logs
				if len(output) > 1000 {
					fmt.Fprintf(debugFile, "Output (first 1000 chars): %s\n", output[:1000])
				} else {
					fmt.Fprintf(debugFile, "Output: %s\n", output)
				}
			} else {
				fmt.Fprintf(debugFile, "Error output: %s\n", errBuf.String())
			}
		}
		
		if err != nil {
			return SubtitleLangMsg{Error: fmt.Sprintf("Error fetching subtitle languages: %v. Check output.log for details.", err)}
		}

		languages := ParseSubtitleLanguages(outBuf.String())
		
		if debugFile != nil {
			fmt.Fprintf(debugFile, "Parsed %d subtitle languages\n", len(languages))
			for i, l := range languages {
				fmt.Fprintf(debugFile, "Language %d: Code=%s, Name=%s\n", i, l.Code, l.Name)
			}
		}
		
		return SubtitleLangMsg{URL: url, Languages: languages}
	}
}

func ParseSubtitleLanguages(output string) []SubtitleLanguage {
	lines := strings.Split(output, "\n")
	var languages []SubtitleLanguage

	// Flag to indicate when we're in the subtitle languages section
	parsingSubtitles := false

	// Parse subtitle language lines
	for _, line := range lines {
		// Check if we've reached the subtitle languages section
		if strings.Contains(line, "Available automatic captions for") {
			parsingSubtitles = true
			continue
		}

		// Skip header lines and empty lines
		if !parsingSubtitles ||
			strings.Contains(line, "Language Name") ||
			strings.Contains(line, "----") ||
			strings.TrimSpace(line) == "" ||
			strings.Contains(line, "[youtube]") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 1 {
			// Extract language code (first field)
			code := fields[0]

			// Skip "en-orig" as it's not a standard language code
			if code == "en-orig" {
				continue
			}

			// Try to get a human-readable name
			name := code
			switch code {
			case "en":
				name = "English"
			case "es":
				name = "Spanish"
			case "fr":
				name = "French"
			case "de":
				name = "German"
			case "it":
				name = "Italian"
			case "pt":
				name = "Portuguese"
			case "ru":
				name = "Russian"
			case "ja":
				name = "Japanese"
			case "ko":
				name = "Korean"
			case "zh-Hans":
				name = "Chinese (Simplified)"
			case "zh-Hant":
				name = "Chinese (Traditional)"
			case "ar":
				name = "Arabic"
			case "hi":
				name = "Hindi"
			case "tr":
				name = "Turkish"
			default:
				// Capitalize first letter for unknown languages
				if len(name) > 0 {
					name = strings.ToUpper(name[:1]) + name[1:]
				}
			}

			lang := SubtitleLanguage{
				Code: code,
				Name: name,
			}

			// Avoid duplicates
			exists := false
			for _, l := range languages {
				if l.Code == lang.Code {
					exists = true
					break
				}
			}

			if !exists {
				languages = append(languages, lang)
			}
		}
	}

	// If no languages found, create a default option
	if len(languages) == 0 {
		languages = append(languages, SubtitleLanguage{
			Code: "en",
			Name: "English",
		})
	}

	return languages
}

func (m AppModel) downloadSubtitles(url string, langCode string) tea.Cmd {
	return func() tea.Msg {
		// Create assets directory if it doesn't exist
		os.MkdirAll("assets", 0755)

		var path, ffmpegPath string
		if isWindows() {
			path = "bin/yt-dlp.exe"
			ffmpegPath = "bin/ffmpeg.exe"
		} else {
			path = "bin/yt-dlp"
			ffmpegPath = "bin/ffmpeg"
		}

		// Create or open log file
		logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return SubtitleDownloadMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		// Create buffers to capture output
		var outBuf, errBuf strings.Builder
		
		// Check if ffmpeg exists
		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil

		// Use better parameters for subtitle download
		var args []string
		args = append(args, "--write-sub", "--write-auto-sub", "--sub-lang", langCode, "--skip-download")
		
		// Add ffmpeg location if it exists
		if useFfmpeg {
			args = append(args, "--ffmpeg-location", ffmpegPath)
		}
		
		// Add rate limiting parameters
		args = append(args, "--sleep-requests", "1", "--sleep-interval", "5", "--max-sleep-interval", "10")
		args = append(args, "-o", "assets/subtitles.%(ext)s", url)
		
		cmd := exec.Command(path, args...)
		cmd.Stdout = io.MultiWriter(&outBuf, logFile)
		cmd.Stderr = io.MultiWriter(&errBuf, logFile)
		err = cmd.Run()
		
		if err != nil {
			errorOutput := errBuf.String()
			// Check if it's a rate limiting error
			if strings.Contains(errorOutput, "429") || strings.Contains(errorOutput, "Too Many Requests") {
				return SubtitleDownloadMsg{Error: "Rate limited by YouTube. Please try again later."}
			}
			return SubtitleDownloadMsg{Error: fmt.Sprintf("Error downloading subtitles: %v. Check output.log for details.", err)}
		}

		return SubtitleDownloadMsg{Done: true}
	}
}

type SubtitleLangMsg struct {
	URL       string
	Languages []SubtitleLanguage
	Error     string
}

type SubtitleDownloadMsg struct {
	Done  bool
	Error string
}
