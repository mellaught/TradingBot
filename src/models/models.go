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
	// If current number of users == NaxMembers the bot willn't continue the dialog.
	NaxMembers int
}

// Struct for members.json
type Members struct {
	M []*User `json:"Users"`
}

// Struct for User.
type User struct {
	ChatId       int64 `json:"ChatId"`
	Notification bool  `json:"Notificaions"`
}
