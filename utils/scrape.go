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
		for i, image := range images {
			imgType := strings.Split(image, ".")[len(strings.Split(image, "."))-1]
			// get seque
			sequence := fmt.Sprintf(".%s", func() string {
				lengthAsString := strconv.Itoa(len(images))
				return strings.Repeat("0", len(strings.Split(lengthAsString, ""))-len(strings.Split(strconv.Itoa(i), ""))) + strconv.Itoa(i)
			}())
			//
			fileName := "./assets/" + pageUrl.Hostname() + sequence + "." + imgType
			// print
			err := saveImage(image, fileName)
			if err != nil {
				continue
			}
		}
		return types.StatusMsg("done")
	}
}

func saveImage(url string, fileName string) error {
	reponse, err := http.Get(url)
	if err != nil {
		return err
	}
	defer reponse.Body.Close()
	// image extention
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reponse.Body)
	if err != nil {
		return err
	}
	return nil
}

func getImages(urlstr string, res *[]string) chromedp.Tasks {
	var images []*cdp.Node

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`:root`),
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
