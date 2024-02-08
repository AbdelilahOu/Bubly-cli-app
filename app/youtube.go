package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// The second view, after a task has been Youtube
func YoutubeView(m AppModel) string {
	c := m.Choice

	tpl := TitleStyle("What youtube tools do you wanna use?") + "\n\n%s"

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("Youtube vedio translator", c == 0),
		checkbox("Youtube vedio downloader", c == 1),
		checkbox("Youtube vedio infos", c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}

// Update loop for the second view after a choice has been made
func UpdateYoutube(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	YoutubeChoices := []string{"carrot planting", "market trip", "reading time"}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(YoutubeChoices) > m.Choice+1 {
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
