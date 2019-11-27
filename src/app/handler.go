package app

import (
	FloydWarshall "TradingBot/src/app/strategy/Floyd-Warshall"
	"log"
	"time"
)

var (
	RunStrategy  chan ExchangeStrategy
	StopStrategy chan ExchangeStrategy
)

// ExchangeStrategy struct for turn on or turn off strategy
type ExchangeStrategy struct {
	Name     string
	Strategy string
}

// StrategyHandler run or stop strategy when user send callback from telegram bot
func (a *App) StrategyHandler() {
	for {
		select {
		case run := <-RunStrategy:
			log.Printf("Started %s strategy for %s Exchange", run.Name, run.Strategy)
			time.Sleep(1 * time.Second)

		case stop := <-StopStrategy:
			log.Printf("Stoped %s strategy for %s Exchange", stop.Name, stop.Strategy)
			time.Sleep(1 * time.Second)
		}
	}
}

func (a *App) StartStrategy(exchange string, strategy string) error {

	switch strategy {
	case "FW":
		switch exchange {
		case "Binance":
			FloydWarshall.Start()
		case "Poloniex":
			
		}
	case "MM":
		switch exchange {
		case "Binance":

		case "Ploniex":
		}

	case "Scalp":
		switch exchange {
		case "Binance":

		case "Ploniex":
		}
	}

	return nil
}
