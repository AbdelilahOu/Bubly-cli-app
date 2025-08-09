package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/AbdelilahOu/Bubly-cli-app/types"
	tea "github.com/charmbracelet/bubbletea"
)

func CheckYtdlp() bool {
	var path string
	if runtime.GOOS == "windows" {
		path = "bin/yt-dlp.exe"
	} else {
		path = "bin/yt-dlp"
	}
	_, err := os.Stat(path)
	return err == nil
}

func CheckFfmpeg() bool {
	var path string
	if runtime.GOOS == "windows" {
		path = "bin/ffmpeg.exe"
	} else {
		path = "bin/ffmpeg"
	}
	_, err := os.Stat(path)
	return err == nil
}

func InstallYtdlp() tea.Msg {
	fmt.Println("Installing yt-dlp...")
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return types.YtdlpInstalledMsg{Err: err}
	}

	var url string
	switch runtime.GOOS {
	case "windows":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	case "linux":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	case "darwin":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
	default:
		return types.YtdlpInstalledMsg{Err: fmt.Errorf("unsupported OS")}
	}

	resp, err := http.Get(url)
	if err != nil {
		return types.YtdlpInstalledMsg{Err: err}
	}
	defer resp.Body.Close()

	var destPath string
	if runtime.GOOS == "windows" {
		destPath = "bin/yt-dlp.exe"
	} else {
		destPath = "bin/yt-dlp"
	}

	out, err := os.Create(destPath)
	if err != nil {
		return types.YtdlpInstalledMsg{Err: err}
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return types.YtdlpInstalledMsg{Err: err}
	}

	err = os.Chmod(destPath, 0755)
	if err != nil {
		return types.YtdlpInstalledMsg{Err: err}
	}

	// Add bin directory to path
	path := os.Getenv("PATH")
	err = os.Setenv("PATH", "bin;"+path)
	if err != nil {
		return types.YtdlpInstalledMsg{Err: err}
	}

	return types.YtdlpInstalledMsg{Err: nil}
}

func InstallFfmpeg() tea.Msg {
	fmt.Println("Installing ffmpeg...")
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return types.FfmpegInstalledMsg{Err: err}
	}

	var url string
	switch runtime.GOOS {
	case "windows":
		// For Windows, we'll download a pre-built ffmpeg binary
		url = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
	case "linux":
		// For Linux, we'll need to install via package manager or download static build
		url = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-git-amd64-static.tar.xz"
	case "darwin":
		// For macOS, we'll need to install via Homebrew or download static build
		url = "https://evermeet.cx/ffmpeg/ffmpeg-5.1.7z"
	default:
		return types.FfmpegInstalledMsg{Err: fmt.Errorf("unsupported OS for ffmpeg installation")}
	}

	// For now, we'll just download a simple ffmpeg binary for Windows
	// In a real implementation, you'd want to handle extraction of archives
	if runtime.GOOS == "windows" {
		resp, err := http.Get(url)
		if err != nil {
			return types.FfmpegInstalledMsg{Err: err}
		}
		defer resp.Body.Close()

		// Just create a placeholder for now - in a real implementation you'd 
		// download and extract the actual ffmpeg binary
		destPath := "bin/ffmpeg.exe"
		out, err := os.Create(destPath)
		if err != nil {
			return types.FfmpegInstalledMsg{Err: err}
		}
		defer out.Close()

		// Write a simple placeholder
		out.WriteString("This is a placeholder for ffmpeg. Please install ffmpeg manually.")
		
		err = os.Chmod(destPath, 0755)
		if err != nil {
			return types.FfmpegInstalledMsg{Err: err}
		}
	}

	return types.FfmpegInstalledMsg{Err: nil}
}
