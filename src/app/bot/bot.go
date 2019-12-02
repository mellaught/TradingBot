package bot

import (
	"TradingBot/src/models"
	"fmt"
	"log"
	"strings"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

//InitBot initialization: loading the database, creating a bot by token from the config.
func InitBot(config *models.BotConfig, members *models.Members) *Bot {

	b := Bot{
		Token:           config.Token,
		Dlg:             map[int64]*Dialog{},
		UserStrategy:    map[int64]string{},
		Members:         map[int64]bool{},
		MembersStrategy: make(map[int64]map[string]*Strategy),
		RunStrategy:     make(chan ExchangeStrategy),
		StopStrategy:    make(chan ExchangeStrategy),
		pass:            config.Password,
	}

	// Create new bot
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Fatal(err)
	}

	b.Bot = bot
	for _, m := range members.M {
		b.Members[m.ChatId] = m.Notification
	}

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
			b.SendMessage(errPassMessage, ChatId, nil)
			return
		} else {
			b.Members[ChatId] = true
			kb := b.MainKb()
			UserHistory[ChatId] = ""
			b.SendMessage(welcomeMessage, ChatId, kb)
			b.WriteToJson(ChatId, true)
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
		if strings.Contains(update.CallbackQuery.Data, yesStrategy) {
			dialog.Text = update.CallbackQuery.Data[11:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:11]
		} else if strings.Contains(update.CallbackQuery.Data, noStrategy) {
			dialog.Text = update.CallbackQuery.Data[10:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:10]
		}
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
	// Check ChatId in b.Members.
	if command == startCommand {
		// "/Start" interacting with the bot, bot description and available commands.
		UserHistory[ChatId] = "start"
		b.SendMessage(startMessage, ChatId, nil)
		return
	}

	if ok, _ := b.Members[ChatId]; !ok {
		UserHistory[ChatId] = "start"
		b.SendMessage(notAuto, ChatId, nil)
		return
	}

	switch command {
	// Get Main Menu Keyboard
	case getMainMenu:
		kb, txt := b.GetMenuMessage(ChatId)
		b.SendMessage(txt, ChatId, kb)
		return

	// Cancel command return prev step for user
	case cancelComm:
		b.CancelHandler(ChatId)

	// Get Trading choose kb
	case tradingCommand:
		UserHistory[ChatId] = "strategies"
		kb := b.StrategiesKb()
		b.EditAndSend(&kb, strategiesMessage, ChatId)
		return

	// Turn ON current strategy
	case yesStrategy:
		b.StrategyHandler(b.UserStrategy[ChatId], b.Dlg[ChatId].Text, true, ChatId)
		txt := fmt.Sprintf("Strategy %s *ON*", b.Dlg[ChatId].Text)
		b.EditAndSend(nil, txt, ChatId)
		b.WriteToJson(ChatId, true)
		return

	// Turn OFF current strategy
	case noStrategy:
		b.StrategyHandler(b.UserStrategy[ChatId], b.Dlg[ChatId].Text, false, ChatId)
		txt := fmt.Sprintf("Strategy %s *OFF*", b.Dlg[ChatId].Text)
		b.EditAndSend(nil, txt, ChatId)
		b.WriteToJson(ChatId, false)
		return

	// Floyd Warshall command
	case fwcommand:
		UserHistory[ChatId] = "FW_0"
		kb := b.YesNoNotifyKb()
		txt := fmt.Sprintf(strategyPower, "Floyd Warshall")
		b.EditAndSend(&kb, txt, ChatId)
		return

	// Get Notify choose kb
	case notifyCommand:
		kb := b.YesNoNotifyKb()
		b.EditAndSend(&kb, notifyMessage, ChatId)
		return

	// Subsctibe notifications
	case yesNotify:
		UserNotifications[ChatId] = true
		b.EditAndSend(nil, "Notifications *ON*", ChatId)
		b.WriteToJson(ChatId, true)
		return

	// Unsubscribe notifications
	case noNotify:
		UserNotifications[ChatId] = false
		b.EditAndSend(nil, "Notifications *OFF*", ChatId)
		b.WriteToJson(ChatId, false)
		return

	// StopBot
	case offBot:
		//b.Bot.StopReceivingUpdates()
	}
}
