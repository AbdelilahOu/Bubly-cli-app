package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var ScrapingOptions = []ViewsOptions{
	{
		View:        "web-images",
		ChoiceLabel: "Download images from web site 🖼️",
	},
	{
		View:        "web-pdf",
		ChoiceLabel: "Print a website 📄",
	},
}

func UpdateScraping(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	if len(m.History) > 2 {
		switch m.History[2] {
		case "web-images":
			return UpdateWebsiteImages(msg, m)
		case "web-pdf":
			return UpdateWebsitePrint(msg, m)
		}
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(ScrapingOptions) > m.Choice+1 {
				m.Choice++
			}
		case "k", "up":
			if m.Choice > 0 {
				m.Choice--
			}
		case "enter":
			m.History = append(m.History, ScrapingOptions[m.Choice].View)
			return m, nil
		}

	}
	return m, nil
}

func ScrapingView(m AppModel) string {
	c := m.Choice

	tpl := TitleStyle("What web scraping tools do you wanna use? 🔨") + "\n\n%s"
	s := " "

	if len(m.History) > 2 {
		switch m.History[2] {
		case "web-images":
			s = WebsiteImagesView(m)
		case "web-pdf":
			s = PrintWebsiteView(m)
		}
	} else {
		s = fmt.Sprintf(tpl, fmt.Sprintf(
			strings.Repeat("%s\n", len(ScrapingOptions)),
			destructureOptions(ScrapingOptions, c)...,
		))
	}

	return s
}

// pdf printer view and update funcs
func PrintWebsiteView(m AppModel) string {
	tpl := TitleStyle("Print a website 📄") + "\n\n"

	return tpl
}

func UpdateWebsitePrint(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// m.History = append(m.History, ScrapingChoices[m.Choice])
			return m, nil
		}

	}
	return m, nil
}

// website images downloader view and update functions
func WebsiteImagesView(m AppModel) string {
	tpl := TitleStyle("Download images from web site 🖼️") + "\n\n"

	return tpl
}

func UpdateWebsiteImages(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// m.History = append(m.History, ScrapingChoices[m.Choice])
			return m, nil
		}

	}
	return m, nil
}
