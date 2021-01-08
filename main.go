package main

import (
	"context"
	"sync"

	"github.com/c-m-hunt/ps5-checker/check"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

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
}

func main() {
	checkers := check.CheckerList{}
	options := check.NewOptions()
	options.Headless = true
	cb := check.CheckerBase{Options: options}

	var ctx context.Context
	cancel := check.SetupBrowserContext(options, &ctx)
	defer cancel()

	checkers = append(checkers, &check.Game{CheckerBase: cb, Context: &ctx})
	checkers = append(checkers, &check.Argos{CheckerBase: cb, Context: &ctx})
	checkers = append(checkers, &check.Smyths{CheckerBase: cb, Context: &ctx})
	checkers = append(checkers, &check.Amazon{CheckerBase: cb, Context: &ctx})

	var wg sync.WaitGroup

	for _, c := range checkers {
		wg.Add(1)
		go runCheck(c, &wg)
	}

	wg.Wait()
}

func runCheck(c check.Checker, wg *sync.WaitGroup) {
	defer wg.Done()
	check.RunStockCheck(c)
}
