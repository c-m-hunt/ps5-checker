package check

import (
	"context"
	"reflect"
	"strings"

	"github.com/chromedp/chromedp"
)

type Smyths struct {
	CheckerBase
}

func (c *Smyths) GetName() string {
	t := reflect.TypeOf(c)
	return t.Elem().Name()
}

func (c Smyths) PrintStatus() {
	c.CheckerInfo.PrintStatus(c.GetName())
}

func (s *Smyths) CheckStock() error {
	s.CheckerInfo.LogCheck()
	url := "https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/playstation-5/playstation-5-consoles/playstation-5-console/p/191259"

	var ctx context.Context
	cancels := SetupBrowserContext(s.Options, &ctx)
	for _, c := range cancels {
		defer c()
	}

	var stock string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitEnabled(".cookieProcessed"),
		chromedp.Click(".cookieProcessed"),
		chromedp.Text(".resultStock", &stock, chromedp.NodeVisible),
		chromedp.WaitVisible(".resultStock"),
	)

	if err != nil {
		s.Errors++
		return err
	}
	if strings.TrimSpace(stock) != "Out Of Stock" {
		s.CheckerInfo.LogStockSeen(s.GetName(), url, ctx)

		return nil
	}

	return nil
}
