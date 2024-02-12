package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AbdelilahOu/Bubly-cli-app/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GetPageAsPdf(URL string) tea.Cmd {
	if _, err := os.Stat("./assets/"); os.IsNotExist(err) {
		_ = os.Mkdir("./assets", 0755)
	}
	return func() tea.Msg {
		// create context
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		// get data
		var buf []byte
		if err := chromedp.Run(ctx, printToPDF(URL, &buf)); err != nil {
			return types.StatusMsg("error")
		}
		// parse base url and use as a name
		pageUrl, err := url.Parse(URL)
		if err != nil {
			return types.StatusMsg("error")
		}
		// file path
		fileName := "./assets/" + pageUrl.Hostname() + ".pdf"
		if err := os.WriteFile(fileName, buf, 0o644); err != nil {
			return types.StatusMsg("error")
		}
		return types.StatusMsg("done")
	}
}

func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`:root`),
		chromedp.Sleep(time.Second * 2),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).WithPaperHeight(12).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func GetPageImages(URL string) tea.Cmd {
	if _, err := os.Stat("./assets/"); os.IsNotExist(err) {
		_ = os.Mkdir("./assets", 0755)
	}
	return func() tea.Msg {
		// create context
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		// get data
		var images []string
		if err := chromedp.Run(ctx, getImages(URL, &images)); err != nil {
			return types.StatusMsg("error")
		}
		pageUrl, err := url.Parse(URL)
		if err != nil {
			return types.StatusMsg("error")
		}
		for i, image := range images[:3] {
			reponse, err := http.Get(image)
			if err != nil {
				fmt.Println()
				continue
			}
			defer reponse.Body.Close()
			// image extention
			imgType := strings.Split(image, ".")[len(strings.Split(image, "."))-1]
			// get seque
			sequence := fmt.Sprintf(".%s", func() string {
				length := len(images)
				lengthAsString := strconv.Itoa(length)
				return strings.Repeat("0", len(strings.Split(lengthAsString, ""))) + strconv.Itoa(i)
			}())
			//
			fileName := "./assets/" + pageUrl.Hostname() + sequence + "." + imgType
			//
			file, err := os.Create(fileName)
			if err != nil {
				continue
			}
			defer file.Close()
			_, err = io.Copy(file, reponse.Body)
			if err != nil {
				continue
			}
		}
		return types.StatusMsg("done")
	}
}

func getImages(urlstr string, res *[]string) chromedp.Tasks {
	var images []*cdp.Node

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`:root`),
		chromedp.Sleep(time.Second * 2),
		chromedp.Nodes("img", &images, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var src string
			var srcs []string
			for _, image := range images {
				src = image.AttributeValue("src")
				srcs = append(srcs, src)
			}
			*res = srcs
			return nil
		}),
	}
}
