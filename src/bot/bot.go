package bot

import (
	"TradingBot/src/models"
	"log"
	"strings"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

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
	Token   string
	Bot     *tgbotapi.BotAPI
	Dlg     map[int64]*Dialog
	Members *models.Members
	pass    string
}

//InitBot initialization: loading the database, creating a bot by token from the config.
func InitBot(config *models.BotConfig, members *models.Members) *Bot {

	b := Bot{
		Token:   config.Token,
		Dlg:     map[int64]*Dialog{},
		Members: &models.Members{},
		pass:    config.Password,
	}

	// Create new bot
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Fatal(err)
	}

	b.Bot = bot
	b.Members = members

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
	if strings.Contains(UserHistory[ChatId], "start") {
		if text != b.pass {
			b.SendMessage("Error password:(\nTry again, my friend!", ChatId, nil)
			return
		} else {
			kb := b.MainKb()
			UserHistory[ChatId] = ""
			b.SendMessage("Welcome to Trading Bot", ChatId, kb)
			return
		}
	}
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
		b.SendMessage(startMessage, ChatId, nil)
		return

	// Get Main Menu Keyboard.
	case getMainMenu:
		kb, txt := b.GetMenuMessage(ChatId)
		b.SendMessage(txt, ChatId, kb)
		return
	
	// Get Notify choose kb.
	case notifyCommand:
		kb := b.YesNoNotifyKb()
		b.EditAndSend(&kb, notifyMessage, ChatId)
		return

	// Subsctibe notifications
	case yesNotify:
		UserNotifications[ChatId] = true
		b.SendMessage("Notificaions ON", ChatId, nil)
		b.WriteToJson(ChatId, true)
		return

	// Unsubscribe notifications
	case noNotify:
		UserNotifications[ChatId] = false
		b.SendMessage("Notificaions OFF", ChatId, nil)
		b.WriteToJson(ChatId, false)
		return

	// StopBot
	case OffBot:
		b.Bot.StopReceivingUpdates()
	}
}
