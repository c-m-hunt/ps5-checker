package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/c-m-hunt/ps5-checker/check"
	"github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
}

func main() {
	checkers := check.CheckerList{}
	options := check.NewOptions()
	options.Headless = true

	checkers = append(checkers, &check.Game{Options: options})
	checkers = append(checkers, &check.Argos{Options: options})
	checkers = append(checkers, &check.Smyths{Options: options})
	checkers = append(checkers, &check.Amazon{Options: options})

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
			beeep.Alert("Stock found", msg, "")
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
