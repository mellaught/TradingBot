package models

//Config ..
type ConfigFile struct {
	Bot *BotConfig
}

// ExchangeConfig holds Api Key and Api Secret for private exchanges Api.
type ExchangeConfig struct {
	ApiKey    string
	ApiSecret string
	Timeout   string
}

// Bot config.
type BotConfig struct {
	// Bot Token
	Token string
	// Password for author to start use bot.
	Password string
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

// Message for User from any strategies
type Message struct {
	ChatId int64
	Txt    string
}
