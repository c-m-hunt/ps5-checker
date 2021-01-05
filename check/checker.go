package check

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type CheckerList []Checker

type Checker interface {
	GetName() string
	CheckStock() error
	GetInStock() bool
	PrintStatus()
}

type CheckerInfo struct {
	inStock       bool
	checks        int
	errors        int
	stockSeen     int
	lastCheck     *time.Time
	stockLastSeen *time.Time
}

type Options struct {
	Headless bool
}

func NewOptions() Options {
	return Options{
		Headless: true,
	}
}

func SetupBrowserContext(o Options) (context.Context, []func()) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", o.Headless),
	)

	ctxTimeout, cancel0 := context.WithTimeout(context.Background(), 20*time.Second)

	allocCtx, cancel1 := chromedp.NewExecAllocator(ctxTimeout, opts...)

	ctx, cancel2 := chromedp.NewContext(allocCtx)

	return ctx, []func(){cancel0, cancel1, cancel2}
}

func (c CheckerInfo) PrintStatus(name string) {
	fmt.Printf("%-10v %8v %8v %8v %-17v %-17v\n", "Name", "Checks", "Errors", "Stock", "Last check", "Stock last seen")
	lastSeen := ""
	if c.stockLastSeen != nil {
		lastSeen = c.stockLastSeen.Format("Jan 2 15:04:05")
	}
	fmt.Printf("%-10v %8v %8v %8v %-17v %-17v\n", name, c.checks, c.errors, c.stockSeen, c.lastCheck.Format("Jan 2 15:04:05"), lastSeen)
}

func (c *CheckerInfo) LogCheck() {
	c.inStock = false
	c.checks++
	now := time.Now()
	c.lastCheck = &now
}

func (c *CheckerInfo) LogStockSeen() {
	c.inStock = true
	now := time.Now()
	c.stockLastSeen = &now
}
