package bot

import (
	"log"
	"tradingBot/src/models"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	CurrentPrice   float64
	CurrnetMarkup  string
	commands       = make(map[int64]string)
	UserHistory    = make(map[int64]string)
	MinterAddress  = make(map[int64]string)
	BitcoinAddress = make(map[int64]string)
	CoinToSell     = make(map[int64]string)
	EmailAddress   = make(map[int64]string)
	PriceToSell    = make(map[int64]float64)
	SaveBuy        = make(map[int64]bool)
	SaveSell       = make(map[int64]bool)
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
	Token string
	Bot   *tgbotapi.BotAPI
	Dlg   map[int64]*Dialog
}

//InitBot initialization: loading the database, creating a bot by token from the config.
func InitBot(config *models.Config) *Bot {

	b := Bot{
		Token: config.BotToken,
		Dlg:   map[int64]*Dialog{},
	}

	// Create new bot
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Fatal(err)
	}

	b.Bot = bot
	go b.Run()

	return &b
}

// Run is starting bot.
func (b *Bot) Run() {

	//Set update timeout
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Get updates from bot
	updates, _ := b.Bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		dialog, exist := b.assembleUpdate(update)
		if !exist {
			continue
		}

		b.Dlg[dialog.ChatId] = dialog

		if botCommand := b.getCommand(update); botCommand != "" {
			b.RunCommand(botCommand, dialog.ChatId)
			continue
		}

		b.TextMessageHandler(update.Message.Text, dialog.ChatId)
		continue
	}
}

// TextMessageHandler
func (b *Bot) TextMessageHandler(text string, ChatId int64) {

}

// assembleUpdate
func (b *Bot) assembleUpdate(update tgbotapi.Update) (*Dialog, bool) {
	dialog := &Dialog{}
	if update.Message != nil {
		dialog.ChatId = update.Message.Chat.ID
		dialog.MessageId = update.Message.MessageID
		dialog.Text = update.Message.Text
	} else if update.CallbackQuery != nil {
		dialog.CallBackId = update.CallbackQuery.ID
		dialog.ChatId = update.CallbackQuery.Message.Chat.ID
		dialog.MessageId = update.CallbackQuery.Message.MessageID
	} else {
		return dialog, false
	}

	command, isset := commands[dialog.ChatId]
	if isset {
		dialog.Command = command
	} else {
		dialog.Command = ""
	}

	return dialog, true
}

// getCommand returns command from telegram update
func (b *Bot) getCommand(update tgbotapi.Update) string {
	if update.Message != nil {
		if update.Message.IsCommand() {
			return update.Message.Command()
		}
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}

	return ""
}

// RunCommand executes the input command.
func (b *Bot) RunCommand(command string, ChatId int64) {
	commands[ChatId] = command
	switch command {
	// "/Start" interacting with the bot, bot description and available commands.
	case startCommand:
		UserHistory[ChatId] = "start"
		b.SendMessage("Hello", ChatId, nil)
		return

	// case cancelComm:
	// 	b.CancelHandler(ChatId)

	case subscribe:
		// Subsctibe
		b.SendMessage("Subscribe succesfull!", ChatId, nil)
		return

	case unsubscribe:
		// Unsubscribe
		b.SendMessage("Unsubscribe succesfull!", ChatId, nil)
		return

	case getMainMenu:
		kb, txt, err := b.SendMenuMessage(ChatId)
		if err != nil {
			b.PrintAndSendError(err, ChatId)
			return
		}

		b.SendMessage(txt, ChatId, kb)
		return
	}
}
