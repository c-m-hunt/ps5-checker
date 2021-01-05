package check

import (
	"reflect"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type Game struct {
	CheckerInfo
	Options
}

func (g *Game) GetName() string {
	t := reflect.TypeOf(g)
	return t.Elem().Name()
}

func (g *Game) GetInStock() bool {
	return g.inStock
}

func (c Game) PrintStatus() {
	c.CheckerInfo.PrintStatus(c.GetName())
}

func (g *Game) CheckStock() error {
	g.inStock = false
	g.checks++
	url := "https://www.game.co.uk/playstation-5"

	ctx, cancels := SetupBrowserContext(g.Options)
	for _, c := range cancels {
		defer c()
	}

	//var res string
	var stockButtons []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Nodes("#contentPanels3 .sectionButton a", &stockButtons, chromedp.NodeVisible),
	)
	if err != nil {
		g.errors++
		return err
	}
	for _, sb := range stockButtons {
		if sb.Children[0].NodeValue != "Out of stock" {
			g.inStock = true
			return nil
		}
	}
	return nil
}
