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
	GetCheckInfo() CheckerInfo
}

type CheckerBase struct {
	CheckerInfo
	Options
}

func (c *CheckerBase) GetInStock() bool {
	return c.InStock
}

func (c CheckerBase) GetCheckInfo() CheckerInfo {
	return c.CheckerInfo
}

type CheckerInfo struct {
	InStock       bool
	Checks        int
	Errors        int
	StockSeen     int
	LastCheck     *time.Time
	StockLastSeen *time.Time
	StockURL      string
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
	if c.StockLastSeen != nil {
		lastSeen = c.StockLastSeen.Format("Jan 2 15:04:05")
	}
	fmt.Printf("%-10v %8v %8v %8v %-17v %-17v\n", name, c.Checks, c.Errors, c.StockSeen, c.LastCheck.Format("Jan 2 15:04:05"), lastSeen)
}

func (c *CheckerInfo) LogCheck() {
	c.InStock = false
	c.Checks++
	now := time.Now()
	c.LastCheck = &now
}

func (c *CheckerInfo) LogStockSeen(url string) {
	c.InStock = true
	now := time.Now()
	c.StockLastSeen = &now
	c.StockURL = url
}
