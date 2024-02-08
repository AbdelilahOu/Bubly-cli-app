package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var ScrapingOptions = []ViewsOptions{
	{
		View:        "images",
		ChoiceLabel: "Download images from web site",
	},
	{
		View:        "pdf",
		ChoiceLabel: "Print a website",
	},
}

func ScrapingView(m AppModel) string {
	c := m.Choice

	tpl := TitleStyle("What web scraping tools do you wanna use?") + "\n\n%s"

	choices := fmt.Sprintf(
		"%s\n%s\n",
		checkbox(ScrapingOptions[0].ChoiceLabel, c == 0),
		checkbox(ScrapingOptions[1].ChoiceLabel, c == 1),
	)

	return fmt.Sprintf(tpl, choices)
}

func UpdateScraping(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
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
			// m.History = append(m.History, ScrapingChoices[m.Choice])
			return m, nil
		}

	}
	return m, nil
}
