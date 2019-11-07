package models

//Config ..
type Config struct {
	Binance  *ApiData
	Bitforex *ApiData
	Yobit    *ApiData
	Bittrex  *ApiData
	LiveCoin *ApiData
	BotToken string
}

// ApiData holds Api Key and Api Secret for private exchanges Api.
type ApiData struct {
	ApiKey    string
	ApiSecret string
}
