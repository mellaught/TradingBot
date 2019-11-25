package main

import (
	"TradingBot/src/app/exchanges/binance"
	"TradingBot/src/config"
	"TradingBot/src/models"
	"fmt"
	"time"
)

func main() {

	cfg := config.NewViperConfig()

	// Create viper struct.
	conf := &models.Config{
		Bot:      &models.BotConfig{},
		Binance:  &models.ExchangeConfig{},
		Bitforex: &models.ExchangeConfig{},
		Bittrex:  &models.ExchangeConfig{},
		Yobit:    &models.ExchangeConfig{},
	}

	// Read config.json -> viper
	conf.Bot = cfg.ReadBot()
	// Binance
	conf.Binance = cfg.ReadExchanges("binance")
	fmt.Println(conf.Binance)
	b := binance.CreateWorker(conf.Binance)
	fmt.Println("App started")
	b.Start()
	//time.Sleep(2 * time.Second)
	b.StrategyStart()
	time.Sleep(30 * time.Minute)
}
