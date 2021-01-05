package check

import (
	"reflect"
	"strings"

	"github.com/chromedp/chromedp"
)

type Smyths struct {
	CheckerInfo
	Options
}

func (g *Smyths) GetName() string {
	t := reflect.TypeOf(g)
	return t.Elem().Name()
}

func (s *Smyths) GetInStock() bool {
	return s.inStock
}

func (c Smyths) PrintStatus() {
	c.CheckerInfo.PrintStatus(c.GetName())
}

func (s *Smyths) CheckStock() error {
	s.CheckerInfo.LogCheck()
	url := "https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/playstation-5/playstation-5-consoles/playstation-5-console/p/191259"

	ctx, cancels := SetupBrowserContext(s.Options)
	for _, c := range cancels {
		defer c()
	}

	var stock string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitEnabled(".cookieProcessed"),
		chromedp.Click(".cookieProcessed"),
		chromedp.Text(".resultStock", &stock, chromedp.NodeVisible),
	)
	if err != nil {
		s.errors++
		return err
	}
	if strings.TrimSpace(stock) != "Out Of Stock" {
		s.CheckerInfo.LogStockSeen()
		return nil
	}

	return nil
}
