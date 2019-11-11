package app

import (
	"TradingBot/src/models"
	"tradingBot/src/app/bot"
)

// Struct for Tranding Bot.
type App struct {
	Bot *bot.Bot
}

// InitService is initializes the app.
func NewApp(conf *models.Config, members *models.Members) *App {

	a := App{
		Bot: &bot.Bot{},
	}

	// Init Bot.
	a.Bot = bot.InitBot(conf.Bot, members)

	return &a
}

func (a *App) Run() {
	go a.Bot.Run()
	// Start Binance Exchange.
	//go a.Exchanges.Binance.Run()
}
