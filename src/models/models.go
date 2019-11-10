package models

//Config ..
type Config struct {
	Binance  *ExchangeConfig
	Bitforex *ExchangeConfig
	Yobit    *ExchangeConfig
	Bittrex  *ExchangeConfig
	LiveCoin *ExchangeConfig
	Bot      *BotConfig
}

// ExchangeConfig holds Api Key and Api Secret for private exchanges Api.
type ExchangeConfig struct {
	ApiKey    string
	ApiSecret string
}

// Bot config.
type BotConfig struct {
	// Bot Token
	Token string
	// Password for author to start use bot.
	Password string
	// Max Members for bot.
	// If current number users == NumberUsers the bot willn't continue the dialog.
	NumberUsers int
}

// Struct for members.json
type Members struct {
	Users []*User
}

// Struct for User.
type User struct {
	ChatId       int64
	Notification bool
}
