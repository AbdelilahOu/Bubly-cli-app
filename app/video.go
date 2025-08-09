package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type VideoFormat struct {
	ID         string
	Format     string
	Quality    string
	Filesize   string
	Resolution string
}

type VideoFormatSelection struct {
	URL         string
	Formats     []VideoFormat
	Choice      int
	Selected    bool
	Downloading bool
	Done        bool
	Error       bool
	ErrMsg      string
}

func (m AppModel) fetchVideoFormats(url string) tea.Cmd {
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
			return VideoFormatMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
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
			fmt.Fprintf(debugFile, "Fetching video formats for URL: %s\n", url)
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
			return VideoFormatMsg{Error: fmt.Sprintf("Error fetching formats: %v. Check output.log for details.", err)}
		}

		formats := ParseVideoFormats(outBuf.String())

		if debugFile != nil {
			fmt.Fprintf(debugFile, "Parsed %d video formats\n", len(formats))
			for i, f := range formats {
				fmt.Fprintf(debugFile, "Video Format %d: ID=%s, Quality=%s\n", i, f.ID, f.Quality)
			}
		}

		return VideoFormatMsg{URL: url, Formats: formats}
	}
}

func ParseVideoFormats(output string) []VideoFormat {
	lines := strings.Split(output, "\n")
	var formats []VideoFormat

	// Parse format lines
	for _, line := range lines {
		// Skip header lines and empty lines
		if strings.Contains(line, "Available formats") ||
			strings.Contains(line, "ID  EXT") ||
			strings.Contains(line, "----") ||
			strings.TrimSpace(line) == "" ||
			strings.Contains(line, "[youtube]") {
			continue
		}

		// Skip audio-only formats for video selection
		if strings.Contains(line, "audio only") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 3 {
			// Extract format information
			id := fields[0]

			// Find format, resolution, and filesize information
			format := "video"
			resolution := "Unknown resolution"
			filesize := "Unknown size"

			// Look for resolution info (e.g., "640x360")
			for _, field := range fields {
				if strings.Contains(field, "x") && strings.ContainsAny(field, "0123456789") {
					resolution = field
				}
				// Look for filesize info (MiB or KiB)
				if strings.Contains(field, "MiB") || strings.Contains(field, "KiB") {
					filesize = field
				}
			}

			// Try to get a better quality description
			quality := resolution
			if strings.Contains(line, "360p") {
				quality = "360p"
			} else if strings.Contains(line, "480p") {
				quality = "480p"
			} else if strings.Contains(line, "720p") {
				quality = "720p HD"
			} else if strings.Contains(line, "1080p") {
				quality = "1080p Full HD"
			} else if strings.Contains(line, "1440p") {
				quality = "1440p Quad HD"
			} else if strings.Contains(line, "2160p") {
				quality = "2160p 4K"
			}

			videoFormat := VideoFormat{
				ID:         id,
				Format:     format,
				Quality:    quality,
				Filesize:   filesize,
				Resolution: resolution,
			}

			// Avoid duplicates
			exists := false
			for _, f := range formats {
				if f.ID == videoFormat.ID {
					exists = true
					break
				}
			}

			if !exists {
				formats = append(formats, videoFormat)
			}
		}
	}

	// If no formats found, create some default options
	if len(formats) == 0 {
		formats = append(formats, VideoFormat{
			ID:         "best",
			Format:     "video",
			Quality:    "Best quality",
			Filesize:   "Unknown size",
			Resolution: "Highest available",
		})
		formats = append(formats, VideoFormat{
			ID:         "worst",
			Format:     "video",
			Quality:    "Low quality",
			Filesize:   "Unknown size",
			Resolution: "Lowest available",
		})
	}

	return formats
}

func (m AppModel) downloadVideo(url string, formatID string) tea.Cmd {
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
			return VideoDownloadMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		// Create buffers to capture output
		var outBuf, errBuf strings.Builder
		
		// Check if ffmpeg exists
		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil

		// Use better parameters for video download
		var args []string
		args = append(args, "-f", formatID)
		
		// Add ffmpeg location if it exists
		if useFfmpeg {
			args = append(args, "--ffmpeg-location", ffmpegPath)
		}
		
		// Add rate limiting parameters
		args = append(args, "--sleep-requests", "1", "--sleep-interval", "5", "--max-sleep-interval", "10")
		args = append(args, "-o", "assets/video.%(ext)s", url)
		
		cmd := exec.Command(path, args...)
		cmd.Stdout = io.MultiWriter(&outBuf, logFile)
		cmd.Stderr = io.MultiWriter(&errBuf, logFile)
		err = cmd.Run()
		
		if err != nil {
			errorOutput := errBuf.String()
			// Check for specific errors
			if strings.Contains(errorOutput, "403") || strings.Contains(errorOutput, "Forbidden") {
				// Try with a different approach - use best if specific format fails
				args = []string{"-f", "best"}
				
				// Add ffmpeg location if it exists
				if useFfmpeg {
					args = append(args, "--ffmpeg-location", ffmpegPath)
				}
				
				// Add rate limiting parameters
				args = append(args, "--sleep-requests", "1", "--sleep-interval", "5", "--max-sleep-interval", "10")
				args = append(args, "-o", "assets/video.%(ext)s", url)
				
				cmd = exec.Command(path, args...)
				cmd.Stdout = io.MultiWriter(&outBuf, logFile)
				cmd.Stderr = io.MultiWriter(&errBuf, logFile)
				err = cmd.Run()
				
				if err != nil {
					return VideoDownloadMsg{Error: fmt.Sprintf("Error downloading video: %v. Check output.log for details.", err)}
				}
			} else {
				return VideoDownloadMsg{Error: fmt.Sprintf("Error downloading video: %v. Check output.log for details.", err)}
			}
		}

		return VideoDownloadMsg{Done: true}
	}
}

type VideoFormatMsg struct {
	URL     string
	Formats []VideoFormat
	Error   string
}

type VideoDownloadMsg struct {
	Done  bool
	Error string
}
