package check

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type Amazon struct {
	inStock  bool
	Headless bool
}

func (a *Amazon) GetName() string {
	return "Amazon"
}

func (a *Amazon) GetInStock() bool {
	return a.inStock
}

func (a *Amazon) CheckStock() error {
	a.inStock = false
	log.Info("Checking Amazon")
	urls := []string{
		"https://www.amazon.co.uk/PlayStation-9395003-5-Console/dp/B08H97NYGP/ref=sr_1_1?dchild=1&keywords=playstation%2B5&qid=1609854382&sr=8-1&th=1",
		"https://www.amazon.co.uk/PlayStation-9395003-5-Console/dp/B08H95Y452/ref=sr_1_1?dchild=1&keywords=playstation%2B5&qid=1609854382&sr=8-1&th=1"}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", a.Headless),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var stock string

	err := chromedp.Run(ctx,
		chromedp.Navigate(urls[0]),
		chromedp.Click("#sp-cc-accept", chromedp.NodeVisible),
	)
	if err != nil {
		fmt.Print(err)
		return err
	}
	for _, u := range urls {
		err := chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.Text("#availability span", &stock),
		)
		if err != nil {
			fmt.Print(err)
			return err
		}

		if strings.TrimSpace(stock) != "Currently unavailable." {
			a.inStock = true
			return nil
		}
	}
	return nil
}
