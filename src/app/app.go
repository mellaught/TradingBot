package app

import (
	"tradingBot/src/bot"
)


// Struct for Tranding Bot.
type App struct {
	Bot        *bot.Bot
	Binance    *bi  
}

// InitService is initializes the app.
func NewApp(conf *models.Config, members models.Members) *App {

	a := App{
		Bot: *bot.Bot{}
	}

	// Start Bot.
	a.Bot = bot.InitBot()
	go a.Bot.Run()
	// Start Binance Exchange.
	go a.Exchanges.Binance.Run()

	return &a
}

