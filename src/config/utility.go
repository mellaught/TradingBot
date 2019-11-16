package config

import "TradingBot/src/models"

func (v *viperConfig) ReadExchanges(exchange string) *models.ExchangeConfig {

	return &models.ExchangeConfig{
		ApiKey:    v.GetString(exchange + ".ApiKey"),
		ApiSecret: v.GetString(exchange + ".ApiSecret"),
		Timeout:   v.GetString(exchange + ".timeout"),
	}

}

func (v *viperConfig) ReadBot() *models.BotConfig {
	return &models.BotConfig{
		Token:    v.GetString("bot.token"),
		Password: v.GetString("bot.password"),
	}
}
