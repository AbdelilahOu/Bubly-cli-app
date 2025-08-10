package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

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
	progress := 0
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		progress += 1
		if progress > 100 {

			err := doInstallYtdlp()
			return types.YtdlpInstalledMsg{Err: err}
		}
		return types.ProgressMsg{
			Progress: progress,
			Total:    100,
			Message:  fmt.Sprintf("Installing yt-dlp... %d%%", progress),
		}
	})
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

func InstallFfmpeg() tea.Cmd {
	progress := 0
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		progress += 1
		if progress > 100 {

			err := doInstallFfmpeg()
			return types.FfmpegInstalledMsg{Err: err}
		}
		return types.ProgressMsg{
			Progress: progress,
			Total:    100,
			Message:  fmt.Sprintf("Installing ffmpeg... %d%%", progress),
		}
	})
}

func doInstallFfmpeg() error {
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return err
	}

	var url string
	switch runtime.GOOS {
	case "windows":

		url = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
	case "linux":

		url = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-git-amd64-static.tar.xz"
	case "darwin":

		url = "https://evermeet.cx/ffmpeg/ffmpeg-5.1.7z"
	default:
		return fmt.Errorf("unsupported OS for ffmpeg installation")
	}

	if runtime.GOOS == "windows" {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		destPath := "bin/ffmpeg.exe"
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
	}

	return nil
}
