package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var YoutubeOptions = []ViewsOptions{
	{
		View:        "yt-translate",
		ChoiceLabel: "Youtube vedio translator 📝",
	},
	{
		View:        "yt-download",
		ChoiceLabel: "Youtube vedio downloader 📥",
	},
	{
		View:        "yt-infos",
		ChoiceLabel: "Youtube vedio infos 📎",
	},
}

func UpdateYoutube(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
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
			// m.History = append(m.History, YoutubeChoices[m.Choice])
			return m, nil
		}

	}
	return m, nil
}

func YoutubeView(m AppModel) string {
	c := m.Choice

	tpl := TitleStyle("What youtube tools do you wanna use? 🔨") + "\n\n%s"

	choices := fmt.Sprintf(
		strings.Repeat("%s\n", len(YoutubeOptions)),
		checkbox(YoutubeOptions[0].ChoiceLabel, c == 0),
		checkbox(YoutubeOptions[1].ChoiceLabel, c == 1),
		checkbox(YoutubeOptions[2].ChoiceLabel, c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}
