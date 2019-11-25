package main

import (
	"TradingBot/src/app/exchanges/poloniex"
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
	// Poloniex
	conf.Poloniex = cfg.ReadExchanges("poloniex")
	fmt.Println(conf.Poloniex)
	//symbols = ["USDC_BTC", "USDC_ETH", "USDC_BTC"]
	p := poloniex.CreateWorker(conf.Poloniex)
	fmt.Println("App started")
	p.Start()
	p.StrategyStart()
	time.Sleep(30 * time.Minute)
}
