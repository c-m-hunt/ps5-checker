package check

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Smyths struct {
	CheckerBase
	Context *context.Context
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

	ctx, cancelTab := chromedp.NewContext(*s.Context)
	defer cancelTab()
	ctx, cancelTO := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTO()

	var stock string

	err := chromedp.Run(ctx, chromedp.Navigate(url))

	if s.CheckerInfo.Checks == 1 {
		err = chromedp.Run(ctx,
			chromedp.WaitEnabled(".cookieProcessed"),
			chromedp.Click(".cookieProcessed"),
		)
	}

	err = chromedp.Run(ctx,

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
