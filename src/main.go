package main

import (
	"TradingBot/src/app"
	"TradingBot/src/models"
	"Tradingbot/src/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	conf.Bot.Token = cfg.GetString("bot.token")
	conf.Bot.Password = cfg.GetString("bot.password")
	conf.Binance.ApiKey = cfg.GetString("binance.ApiKey")
	conf.Binance.ApiSecret = cfg.GetString("binance.ApiSecret")

	// Read file with bot's members.
	members := &models.Members{}
	jsonFile, err := os.Open("members.json")
	if err != nil {
		if err.Error() == "open members.json: The system cannot find the file specified." {
			log.Println("First use!")
		} else {
			log.Fatal(err)
		}
	} else {
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &members)
	}

	fmt.Println(members)
	app := app.NewApp(conf, members)
	// Start App
	app.Run()
	log.Println("App started!")
	//time.Sleep(15 * time.Minute)
}
