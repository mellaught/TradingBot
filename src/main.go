package main

import (
	"fmt"
	"log"
	"time"
	"tradingBot/src/config"
	"tradingBot/src/exchanges/binance"
	"tradingBot/src/models"
)

func main() {

	cfg := config.NewViperConfig()

	conf := &models.Config{
		BotToken: cfg.GetString("bot.token"),
		Binance:  &models.ApiData{},
		Bitforex: &models.ApiData{},
		Bittrex:  &models.ApiData{},
		Yobit:    &models.ApiData{},
	}

	conf.Binance.ApiKey = cfg.GetString("binance.ApiKey")
	conf.Binance.ApiSecret = cfg.GetString("binance.ApiSecret")

	b := binance.CreateBinanceWrapper(conf.Binance)

	var dates = []string{"2019-11-05T12:46:26.0Z", "2019-11-05T12:49:26.0Z", "2019-11-05T12:50:26.0Z", "2019-11-05T12:53:26.0Z", "2019-11-05T12:55:26.0Z", "2019-11-05T12:58:26.0Z",
		"2019-11-05T12:59:26.0Z", "2019-11-05T13:02:26.0Z", "2019-11-05T13:03:26.0Z", "2019-11-05T13:06:26.0Z"}

	for i := 1; i < len(dates); i = i + 2 {
		fmt.Println(dates[i-1], dates[i])
		start, err := time.Parse(time.RFC3339, dates[i-1])
		if err != nil {
			log.Println("Error while parsing date start:", err)
		}
		end, err := time.Parse(time.RFC3339, dates[i])
		if err != nil {
			log.Println("Error while parsing date end:", err)
		}
		startTime := start.UnixNano() / 1e6
		endTime := end.UnixNano() / 1e6
		b.GetHistoryTrades(startTime, endTime, i+1)
	}
}
