package app

import (
	"TradingBot/src/app/exchanges/binance"
	"TradingBot/src/app/exchanges/poloniex"
	FloydWarshall "TradingBot/src/app/strategy/Floyd-Warshall"
	"TradingBot/src/models"
	"log"
	"tradingBot/src/app/bot"
)

// Struct for Tranding Bot.
type App struct {
	Bot       *bot.Bot
	Binance   *binance.BinanceWorker
	Poloniex  *poloniex.PoloniexWorker
	FW        *FloydWarshall.FloydWll
	Exchanges map[string]bool
}

// InitService is initializes the app.
func NewApp(conf *models.ConfigFile, members *models.Members, exchanges map[string]bool) *App {

	a := App{
		Bot:       &bot.Bot{},
		Binance:   &binance.BinanceWorker{},
		Poloniex:  &poloniex.PoloniexWorker{},
		FW:        &FloydWarshall.FloydWll{},
		Exchanges: exchanges,
	}

	// Create Bot
	a.Bot = bot.InitBot(conf.Bot, members)

	// Create all Exchanges
	for k, v := range exchanges {
		if v {
			a.CreateExchange(conf.Binance, k)
		}
	}

	// Create FloyWarhall Arbitrage Strategy
	a.FW = FloydWarshall.CreateFloydWarshall()

	return &a
}

// Create current exchange by input name
func (a *App) CreateExchange(conf *models.ExchangeConfig, exchange string) {
	// Exchange's name
	switch exchange {
	case "Binance":
		a.Binance = binance.CreateWorker(conf)
		log.Println("Binance created")

	case "Poloniex":
		a.Poloniex = poloniex.CreateWorker(conf)
		log.Println("Poloniex created")
	}

}

// Start current exchange by input name
func (a *App) StartExchage(exchange string) error {

	return nil
}

// Run BOT and Exchanges
func (a *App) Run() {
	// Start Bot
	a.Bot.Run()
	// Start Strategy Handler
	a.StrategyHandler()
	// Start all Exchanges
	for k, v := range a.Exchanges {
		if v {
			err := a.StartExchage(k)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
