package app

import (
	"TradingBot/src/app/exchanges/binance"
	"TradingBot/src/app/exchanges/poloniex"
	FloydWarshall "TradingBot/src/app/strategy/Floyd-Warshall"
	"TradingBot/src/models"
	"fmt"
	"log"
	"tradingBot/src/app/bot"
)

// Struct for Tranding Bot.
type App struct {
	Bot       *bot.Bot
	Binance   *binance.BinanceWorker
	Poloniex  *poloniex.PoloniexWorker
	FW        *FloydWarshall.FloydWll
	Exchanges map[string]*models.ExchangeConfig
}

// InitService is initializes the app.
func NewApp(conf *models.ConfigFile, members *models.Members, exchanges map[string]*models.ExchangeConfig) *App {

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
		a.CreateExchange(v, k)
	}

	// Create FloyWarhall Arbitrage Strategy
	a.FW = FloydWarshall.CreateFloydWarshall()
	log.Println("FloydWarshall created")

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
func (a *App) StartExchage(exchange string) {
	// Exchange's name
	switch exchange {
	case "Binance":
		go a.Binance.Start()
		log.Println("Binance started")

	case "Poloniex":
		go a.Poloniex.Start()
		log.Println("Poloniex started")
	}

}

// Run BOT and Exchanges
func (a *App) Run() {
	fmt.Println("In RUN()")
	// Start Bot
	go a.Bot.Run()
	log.Println("Telegram bot started!")
	// Start Strategy Handler
	go a.StrategyHandler()
	log.Println("Strategy handler started!")
	// Start all Exchanges
	for k, _ := range a.Exchanges {
		a.StartExchage(k)
	}
}
