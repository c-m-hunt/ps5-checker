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
	inStock   bool
	checks    int
	errors    int
	lastCheck time.Time
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
	fmt.Printf("%-10v %8v %8v\n", "Name", "Checks", "Errors")
	fmt.Printf("%-10v %8v %8v\n", name, c.checks, c.errors)
}
