package check

import (
	"context"
	"strings"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type Smyths struct {
	inStock  bool
	Headless bool
}

func (s *Smyths) GetName() string {
	return "Smyths"
}

func (s *Smyths) GetInStock() bool {
	return s.inStock
}

func (s *Smyths) CheckStock() error {
	s.inStock = false
	log.Info("Checking Smyths")
	url := "https://www.smythstoys.com/uk/en-gb/video-games-and-tablets/playstation-5/playstation-5-consoles/playstation-5-console/p/191259"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", s.Headless),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var stock string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitEnabled(".cookieProcessed"),
		chromedp.Click(".cookieProcessed"),
		chromedp.Text(".resultStock", &stock, chromedp.NodeVisible),
	)
	if err != nil {
		return err
	}
	if strings.TrimSpace(stock) != "Out Of Stock" {
		s.inStock = true
		return nil
	}

	return nil
}
