package main

import (
	"TradingBot/src/app"
	"TradingBot/src/config"
	"TradingBot/src/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	cfg := config.NewViperConfig()

	// Create viper struct.
	conf := &models.Config{
		BotToken: &models.BotConfig{},
		Binance:  &models.ApiData{},
		Bitforex: &models.ApiData{},
		Bittrex:  &models.ApiData{},
		Yobit:    &models.ApiData{},
	}

	// Read config.json -> viper
	conf.Bot.Token = cfg.GetString("bot.token")
	conf.Bot.Password = cfg.GetString("bot.password")
	conf.Bot.MaxMembers = cfg.GetString("bot.Nmembers")
	conf.Binance.ApiKey = cfg.GetString("binance.ApiKey")
	conf.Binance.ApiSecret = cfg.GetString("binance.ApiSecret")

	// Read file with bot's members.
	jsonFile, err := os.Open("app/arbitrator/html/default.json")
	if err != nil {
		if err.Error() == "open app/arbitrator/html/default.json: The system cannot find the file specified." {
			log.Println("First use")

		} else {
			log.Fatal(err)
		}
	} else {
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &data)
	}

	app := app.NewApp(conf, members)
	// Start App
	app.Run()
	log.Println("App started!")
}
