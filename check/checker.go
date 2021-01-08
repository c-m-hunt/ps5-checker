package check

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
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

func SetupBrowserContext(o Options, ctx *context.Context) func() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", o.Headless),
	)
	//ctxTimeout, _ := context.WithTimeout(context.Background(), 20*time.Second)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctxNew, cancel := chromedp.NewContext(allocCtx)
	if err := chromedp.Run(ctxNew); err != nil {
		panic(err)
	}
	*ctx = ctxNew
	return cancel
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

func (c *CheckerInfo) LogStockSeen(name string, url string, ctx context.Context) {
	var buf []byte
	chromedp.Run(ctx,
		fullScreenshot(80, &buf),
	)
	if err := ioutil.WriteFile(fmt.Sprintf("./screens/%v_ss.png", name), buf, 0o644); err != nil {
		log.Fatal(err)
	}
	c.InStock = true
	now := time.Now()
	c.StockLastSeen = &now
	c.StockURL = url
}

func fullScreenshot(quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
