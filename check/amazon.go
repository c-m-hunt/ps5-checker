package check

import (
	"reflect"
	"strings"

	"github.com/chromedp/chromedp"
)

type Amazon struct {
	inStock bool
	CheckerInfo
	Options
}

func (g *Amazon) GetName() string {
	t := reflect.TypeOf(g)
	return t.Elem().Name()
}

func (a *Amazon) GetInStock() bool {
	return a.inStock
}

func (c Amazon) PrintStatus() {
	c.CheckerInfo.PrintStatus(c.GetName())
}

func (a *Amazon) CheckStock() error {
	a.CheckerInfo.LogCheck()
	urls := []string{
		"https://www.amazon.co.uk/PlayStation-9395003-5-Console/dp/B08H97NYGP/ref=sr_1_1?dchild=1&keywords=playstation%2B5&qid=1609854382&sr=8-1&th=1",
		"https://www.amazon.co.uk/PlayStation-9395003-5-Console/dp/B08H95Y452/ref=sr_1_1?dchild=1&keywords=playstation%2B5&qid=1609854382&sr=8-1&th=1"}

	ctx, cancels := SetupBrowserContext(a.Options)
	for _, c := range cancels {
		defer c()
	}

	var stock string

	err := chromedp.Run(ctx,
		chromedp.Navigate(urls[0]),
		chromedp.Click("#sp-cc-accept", chromedp.NodeVisible),
	)
	if err != nil {
		a.errors++
		return err
	}
	for _, u := range urls {
		err := chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.Text("#availability span", &stock),
		)
		if err != nil {
			a.errors++
			return err
		}

		if strings.TrimSpace(stock) != "Currently unavailable." {
			a.CheckerInfo.LogStockSeen()
			return nil
		}
	}
	return nil
}
