package app

import (
	"tradingBot/src/bot"
)


// App is LightningOracle main app.
type App struct {
	Bot        *bot.Bot
	
}

// InitService is initializes the app.
func NewApp(conf *models.Config) *App {

	a := App{
		Bot: arbitrator.NewArbitratorApp(conf),

	}

	return &a
}

