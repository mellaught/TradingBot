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
			a.StartStrategy(run.Ctx, run.Exchange, run.Strategy, run.ChatId)
			log.Printf("Started %s strategy for %s Exchange", run.Exchange, run.Strategy)
			time.Sleep(1 * time.Second)

		case message := <-a.FW.Messages:
			a.Bot.SendMessage(message.Txt, message.ChatId, nil)
		}
	}
}

// StartStrategy starts strategy
func (a *App) StartStrategy(ctx *context.Context, exchange, strategy string, ChatId int64) {

	switch strategy {
	case "FW":
		switch exchange {
		case "Binance":
			a.FW.Start(ctx, a.Binance, ChatId)
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
}
