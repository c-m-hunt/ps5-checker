package check

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type Game struct {
	inStock  bool
	Headless bool
}

func (g *Game) GetName() string {
	return "Game"
}

func (g *Game) GetInStock() bool {
	return g.inStock
}

func (g *Game) CheckStock() error {
	g.inStock = false
	log.Info("Checking Game")
	url := "https://www.game.co.uk/playstation-5"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", g.Headless),
	)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allocCtx, cancel := chromedp.NewExecAllocator(ctxTimeout, opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	//var res string
	var stockButtons []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Nodes("#contentPanels3 .sectionButton a", &stockButtons, chromedp.NodeVisible),
	)
	if err != nil {
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
