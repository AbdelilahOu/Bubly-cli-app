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

func InstallYtdlp() tea.Cmd {
	return func() tea.Msg {
		err := doInstallYtdlp()
		return types.YtdlpInstalledMsg{Err: err}
	}
}

func doInstallYtdlp() error {
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return err
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
		return fmt.Errorf("unsupported OS")
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
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
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chmod(destPath, 0755)
	if err != nil {
		return err
	}

	path := os.Getenv("PATH")
	err = os.Setenv("PATH", "bin;"+path)
	if err != nil {
		return err
	}

	return nil
}
