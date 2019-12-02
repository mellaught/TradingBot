package bot

import (
	"context"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// Dialog is struct for dialog with user:   - ChatId: User's ChatID
//											- UserId:   Struct App for Rest Api methods
//											- MessageId: Last Message id
//											- Text:   	Text of the last message from the user
//											- language: User's current language
//											- Command: Last command from user
type Dialog struct {
	ChatId     int64
	CallBackId string
	MessageId  int
	Text       string
	Command    string
}

// Bot is struct for Bot:   - Token: secret token from .env
//							- Api:   Struct App for Rest Api methods
//							- DB:    Postgres DB fro users and user's loots.
//							- Bot:	 tgbotapi Bot(token)
//							- Dlg:   For dialog struct
type Bot struct {
	Token           string
	Bot             *tgbotapi.BotAPI
	Dlg             map[int64]*Dialog
	Members         map[int64]bool
	MembersStrategy map[int64]map[string]*Strategy
	RunStrategy     chan ExchangeStrategy
	StopStrategy    chan ExchangeStrategy
	pass            string
}

// Strategy struct for turn on or turn off strategy
type Strategy struct {
	Strategy string
	Ctx      *context.Context
}

// ExchangeStrategy struct for turn on or turn off strategy
type ExchangeStrategy struct {
	Name     string
	Strategy string
	Ctx      *context.Context
}

var (
	commands          = make(map[int64]string)
	UserHistory       = make(map[int64]string)
	UserNotifications = make(map[int64]bool)
	MinterAddress     = make(map[int64]string)
	BitcoinAddress    = make(map[int64]string)
	CoinToSell        = make(map[int64]string)
	EmailAddress      = make(map[int64]string)
	PriceToSell       = make(map[int64]float64)
	SaveBuy           = make(map[int64]bool)
	SaveSell          = make(map[int64]bool)
)
