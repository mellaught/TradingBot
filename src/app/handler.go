package app

import (
	"context"
	"log"
	"time"
)

// StrategyHandler run or stop strategy when user send callback from telegram bot
func (a *App) StrategyHandler() {
	for {
		select {
		case run := <-a.Bot.RunStrategy:
			err := a.StartStrategy(run.Ctx, run.Name, run.Strategy)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Started %s strategy for %s Exchange", run.Name, run.Strategy)
			time.Sleep(1 * time.Second)

		case stop := <-a.Bot.StopStrategy:
			log.Printf("Stoped %s strategy for %s Exchange", stop.Name, stop.Strategy)
			time.Sleep(1 * time.Second)
		}
	}
}

func (a *App) StartStrategy(ctx *context.Context, exchange, strategy string) error {

	switch strategy {
	case "FW":
		switch exchange {
		case "Binance":
			a.FW.Start(ctx, a.Binance)
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
