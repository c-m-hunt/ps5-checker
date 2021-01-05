package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/c-m-hunt/ps5-checker/check"
	"github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	for {
		log.Debug("Checking stock")
		runChecks()
		time.Sleep(10 * time.Second)
	}
}

func runChecks() {
	checkers := check.CheckerList{}

	checkers = append(checkers, &check.Game{Headless: true})
	checkers = append(checkers, &check.Argos{Headless: true})
	checkers = append(checkers, &check.Smyths{Headless: true})
	checkers = append(checkers, &check.Amazon{Headless: true})

	var wg sync.WaitGroup

	for _, c := range checkers {
		wg.Add(1)
		go runCheck(c, &wg)
	}
	wg.Wait()
}

func runCheck(c check.Checker, wg *sync.WaitGroup) {
	defer wg.Done()
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
		log.Debug(fmt.Sprintf("No stock at %v", c.GetName()))
	}
}
