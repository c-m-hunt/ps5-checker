package check

import (
	"context"
	"reflect"

	"github.com/chromedp/chromedp"
)

type Argos struct {
	CheckerBase
}

func (g *Argos) GetName() string {
	t := reflect.TypeOf(g)
	return t.Elem().Name()
}

func (c Argos) PrintStatus() {
	c.CheckerInfo.PrintStatus(c.GetName())
}

func (a *Argos) CheckStock() error {
	a.CheckerInfo.LogCheck()
	urls := []string{"https://www.argos.co.uk/product/8349024", "https://www.argos.co.uk/product/8349000"}

	var ctx context.Context
	cancels := SetupBrowserContext(a.Options, &ctx)
	for _, c := range cancels {
		defer c()
	}

	outOfStockURL := "https://www.argos.co.uk/vp/oos/ps5.html"
	for _, u := range urls {
		var navURL string
		err := chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.WaitReady("body"),
			chromedp.Location(&navURL),
		)
		if err != nil {
			a.Errors++
			return err
		}
		if navURL != outOfStockURL {
			a.CheckerInfo.LogStockSeen(a.GetName(), u, ctx)
			return nil
		}
	}

	return nil
}
