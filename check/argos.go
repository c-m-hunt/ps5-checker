package check

import (
	"reflect"

	"github.com/chromedp/chromedp"
)

type Argos struct {
	CheckerInfo
	Options
}

func (g *Argos) GetName() string {
	t := reflect.TypeOf(g)
	return t.Elem().Name()
}

func (a *Argos) GetInStock() bool {
	return a.inStock
}

func (c Argos) PrintStatus() {
	c.CheckerInfo.PrintStatus(c.GetName())
}

func (a *Argos) CheckStock() error {
	a.checks++
	a.inStock = false
	urls := []string{"https://www.argos.co.uk/product/8349024", "https://www.argos.co.uk/product/8349000"}

	ctx, cancels := SetupBrowserContext(a.Options)
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
			a.errors++
			return err
		}
		if navURL != outOfStockURL {
			a.inStock = true
			return nil
		}
	}

	return nil
}
