package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/c-m-hunt/ps5-checker/check"
	"github.com/gen2brain/beeep"
	"github.com/gregdel/pushover"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var pushoverApp *pushover.Pushover
var pushoverRecipient *pushover.Recipient

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pushoverApp = pushover.New(os.Getenv("PUSHOVER_TOKEN"))
	pushoverRecipient = pushover.NewRecipient(os.Getenv("PUSHOVER_RECIPIENT"))
}

func main() {
	checkers := check.CheckerList{}
	options := check.NewOptions()
	options.Headless = true
	cb := check.CheckerBase{Options: options}

	checkers = append(checkers, &check.Game{CheckerBase: cb})
	checkers = append(checkers, &check.Argos{CheckerBase: cb})
	checkers = append(checkers, &check.Smyths{CheckerBase: cb})
	checkers = append(checkers, &check.Amazon{CheckerBase: cb})

	var wg sync.WaitGroup

	for _, c := range checkers {
		wg.Add(1)
		go runCheck(c, &wg)
	}

	wg.Wait()
}

func runCheck(c check.Checker, wg *sync.WaitGroup) {
	defer wg.Done()
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

func sendAlert(c check.Checker) {
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
