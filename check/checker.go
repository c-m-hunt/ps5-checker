package check

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gen2brain/beeep"
	"github.com/gregdel/pushover"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
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

var pushoverApp *pushover.Pushover
var pushoverRecipient *pushover.Recipient

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pushoverApp = pushover.New(os.Getenv("PUSHOVER_TOKEN"))
	pushoverRecipient = pushover.NewRecipient(os.Getenv("PUSHOVER_RECIPIENT"))
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

func RunStockCheck(c Checker) {
	counter := 0
GetStuckIn:
	for {
		log.Infof("Checking %v", c.GetName())
		err := c.CheckStock()
		if err != nil {
			log.Error(err)
			log.Error(fmt.Sprintf("Problem getting stock from %v", c.GetName()))
		}
		if c.GetInStock() {
			msg := fmt.Sprintf("FOUND STOCK at %v", c.GetName())
			log.Info(msg)
			sendAlert(c)
			beeep.Alert("Stock found", msg, "")
			break GetStuckIn
		} else {
			log.Warn(fmt.Sprintf("No stock at %v", c.GetName()))
		}
		counter++
		if counter%10 == 0 {
			c.PrintStatus()
		}
		time.Sleep(10 * time.Second)
	}
}

func sendAlert(c Checker) {
	message := pushover.NewMessage(fmt.Sprintf("PS5 stock found at %v", c.GetName()))
	message.Title = "PS5 Stock Found"
	message.Priority = pushover.PriorityHigh
	ci := c.GetCheckInfo()
	message.URL = ci.StockURL
	message.URLTitle = c.GetName()

	_, err := pushoverApp.SendMessage(message, pushoverRecipient)
	if err != nil {
		log.Error("Error sending alert")
		log.Error(err)
	}
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
