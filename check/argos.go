package check

import (
	"context"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type Argos struct {
	inStock  bool
	Headless bool
}

func (a *Argos) GetName() string {
	return "Argos"
}

func (a *Argos) GetInStock() bool {
	return a.inStock
}

func (a *Argos) CheckStock() error {
	a.inStock = false
	log.Info("Checking Argos")
	urls := []string{"https://www.argos.co.uk/product/8349024", "https://www.argos.co.uk/product/8349000"}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", a.Headless),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	outOfStockURL := "https://www.argos.co.uk/vp/oos/ps5.html"
	for _, u := range urls {
		var navURL string
		err := chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.WaitReady("body"),
			chromedp.Location(&navURL),
		)
		if err != nil {
			return err
		}
		if navURL != outOfStockURL {
			a.inStock = true
			return nil
		}
	}

	return nil
}
