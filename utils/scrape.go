package utils

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/AbdelilahOu/Bubly-cli-app/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GetPageImages(URL string) tea.Cmd {
	if _, err := os.Stat("./assets/"); os.IsNotExist(err) {
		_ = os.Mkdir("./assets", 0755)
	}
	return func() tea.Msg {
		// create context
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		// get data
		var buf [][]byte
		if err := chromedp.Run(ctx, getImages(URL, &buf)); err != nil {
			return types.StatusMsg("error")
		}
		// parse base url and use as a name
		pageUrl, err := url.Parse(URL)
		if err != nil {
			return types.StatusMsg("error")
		}
		// file path
		fileName := "./assets/" + pageUrl.Hostname() + ".pdf"
		if err := os.WriteFile(fileName, buf[0], 0o644); err != nil {
			return types.StatusMsg("error")
		}
		return types.StatusMsg("done")
	}
}

// print a specific pdf page.
func getImages(urlstr string, res *[][]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`:root`),
		chromedp.Sleep(time.Second * 2),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).WithPaperHeight(12).Do(ctx)
			if err != nil {
				return err
			}
			*res = [][]byte{buf}
			return nil
		}),
	}
}
