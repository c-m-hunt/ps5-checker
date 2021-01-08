package check

import (
	"context"
	"reflect"
	"time"

	"github.com/chromedp/chromedp"
)

type Argos struct {
	CheckerBase
	Context *context.Context
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

	ctx, cancelTab := chromedp.NewContext(*a.Context)
	defer cancelTab()
	ctx, cancelTO := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTO()

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
