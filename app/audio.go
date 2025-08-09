package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type AudioFormat struct {
	ID       string
	Format   string
	Quality  string
	Filesize string
}

type AudioFormatSelection struct {
	URL         string
	Formats     []AudioFormat
	Choice      int
	Selected    bool
	Downloading bool
	Done        bool
	Error       bool
	ErrMsg      string
}

func (m AppModel) fetchAudioFormats(url string) tea.Cmd {
	return func() tea.Msg {
		// Create assets directory if it doesn't exist
		os.MkdirAll("assets", 0755)

		// Run yt-dlp to get format information
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
			return AudioFormatMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		// Create buffers to capture output
		var outBuf, errBuf strings.Builder
		
		// Check if ffmpeg exists
		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil

		// Prepare arguments
		var args []string
		args = append(args, "-F", url)
		
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
			fmt.Fprintf(debugFile, "Fetching formats for URL: %s\n", url)
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
			return AudioFormatMsg{Error: fmt.Sprintf("Error fetching formats: %v. Check output.log for details.", err)}
		}

		formats := ParseAudioFormats(outBuf.String())
		
		if debugFile != nil {
			fmt.Fprintf(debugFile, "Parsed %d formats\n", len(formats))
			for i, f := range formats {
				fmt.Fprintf(debugFile, "Format %d: ID=%s, Quality=%s\n", i, f.ID, f.Quality)
			}
		}
		
		return AudioFormatMsg{URL: url, Formats: formats}
	}
}

func ParseAudioFormats(output string) []AudioFormat {
	lines := strings.Split(output, "\n")
	var formats []AudioFormat

	for _, line := range lines {
		if strings.Contains(line, "Available formats") ||
			strings.Contains(line, "ID  EXT") ||
			strings.Contains(line, "----") ||
			strings.TrimSpace(line) == "" ||
			strings.Contains(line, "[youtube]") {
			continue
		}

		if strings.Contains(line, "audio only") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				id := fields[0]

				// Skip formats with "-drc" as they seem to have issues
				if strings.Contains(id, "-drc") {
					continue
				}

				quality := "Audio"
				filesize := "Unknown size"

				// Look for bitrate information
				for _, field := range fields {
					if strings.HasSuffix(field, "k") {
						bitrateStr := strings.TrimSuffix(field, "k")
						if _, err := fmt.Sscanf(bitrateStr, "%f", new(float64)); err == nil {
							quality = bitrateStr + " kbps"
						}
					}

					if strings.Contains(field, "MiB") || strings.Contains(field, "KiB") {
						filesize = field
					}
				}

				// Try to get better quality descriptions
				if quality == "Audio" {
					if strings.Contains(line, "Default, high") {
						quality = "High quality"
					} else if strings.Contains(line, "Default, low") {
						quality = "Low quality"
					} else if strings.Contains(line, "[en]") {
						quality = "English audio"
					}
				}

				// Determine format type from extension
				formatType := "audio"
				if len(fields) > 1 {
					ext := fields[1]
					if ext == "m4a" {
						formatType = "M4A (AAC)"
					} else if ext == "webm" {
						formatType = "WebM (Opus)"
					}
				}

				format := AudioFormat{
					ID:       id,
					Format:   formatType,
					Quality:  quality,
					Filesize: filesize,
				}

				// Avoid duplicates
				exists := false
				for _, f := range formats {
					if f.ID == format.ID {
						exists = true
						break
					}
				}

				if !exists {
					formats = append(formats, format)
				}
			}
		}
	}

	// Sort formats by quality (higher bitrate first)
	sort.Slice(formats, func(i, j int) bool {
		// Extract bitrate numbers for comparison
		iBitrate := extractBitrate(formats[i].Quality)
		jBitrate := extractBitrate(formats[j].Quality)
		
		// If we can't extract bitrates, sort alphabetically
		if iBitrate == 0 && jBitrate == 0 {
			return formats[i].Quality > formats[j].Quality
		}
		
		return iBitrate > jBitrate
	})

	// If no formats found, create some default options
	if len(formats) == 0 {
		formats = append(formats, AudioFormat{
			ID:       "bestaudio",
			Format:   "audio",
			Quality:  "Best quality",
			Filesize: "Unknown size",
		})
		formats = append(formats, AudioFormat{
			ID:       "worstaudio",
			Format:   "audio",
			Quality:  "Low quality",
			Filesize: "Unknown size",
		})
	}

	return formats
}

// Helper function to extract bitrate from quality string
func extractBitrate(quality string) int {
	re := regexp.MustCompile(`(\d+)\s*kbps`)
	matches := re.FindStringSubmatch(quality)
	if len(matches) > 1 {
		if bitrate, err := strconv.Atoi(matches[1]); err == nil {
			return bitrate
		}
	}
	return 0
}

func (m AppModel) downloadAudio(url string, formatID string) tea.Cmd {
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
			return AudioDownloadMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		// Create buffers to capture output
		var outBuf, errBuf strings.Builder
		
		// Check if ffmpeg exists
		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil
		
		// Use better parameters for audio download
		var args []string
		args = append(args, "-f", formatID, "-x", "--audio-quality", "0")
		
		// Add ffmpeg location if it exists
		if useFfmpeg {
			args = append(args, "--ffmpeg-location", ffmpegPath)
		}
		
		// Add rate limiting parameters
		args = append(args, "--sleep-requests", "1", "--sleep-interval", "5", "--max-sleep-interval", "10")
		args = append(args, "-o", "assets/audio.%(ext)s", url)
		
		cmd := exec.Command(path, args...)
		cmd.Stdout = io.MultiWriter(&outBuf, logFile)
		cmd.Stderr = io.MultiWriter(&errBuf, logFile)
		err = cmd.Run()
		
		if err != nil {
			errorOutput := errBuf.String()
			// Check for specific errors
			if strings.Contains(errorOutput, "403") || strings.Contains(errorOutput, "Forbidden") {
				// Try with a different approach - use bestaudio if specific format fails
				args = []string{"-f", "bestaudio", "-x", "--audio-quality", "0"}
				
				// Add ffmpeg location if it exists
				if useFfmpeg {
					args = append(args, "--ffmpeg-location", ffmpegPath)
				}
				
				// Add rate limiting parameters
				args = append(args, "--sleep-requests", "1", "--sleep-interval", "5", "--max-sleep-interval", "10")
				args = append(args, "-o", "assets/audio.%(ext)s", url)
				
				cmd = exec.Command(path, args...)
				cmd.Stdout = io.MultiWriter(&outBuf, logFile)
				cmd.Stderr = io.MultiWriter(&errBuf, logFile)
				err = cmd.Run()
				
				if err != nil {
					return AudioDownloadMsg{Error: fmt.Sprintf("Error downloading audio: %v. Check output.log for details.", err)}
				}
			} else {
				return AudioDownloadMsg{Error: fmt.Sprintf("Error downloading audio: %v. Check output.log for details.", err)}
			}
		}

		return AudioDownloadMsg{Done: true}
	}
}

type AudioFormatMsg struct {
	URL     string
	Formats []AudioFormat
	Error   string
}

type AudioDownloadMsg struct {
	Done  bool
	Error string
}

func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") ||
		strings.HasSuffix(strings.ToLower(os.Getenv("PATH")), ".exe")
}
