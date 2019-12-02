package bot

import (
	"TradingBot/src/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand   = "start"
	getMainMenu    = "home"
	tradingCommand = "trading"
	tricommand     = "tri"
	fwcommand      = "fw"
	notifyCommand  = "notify"
	settingsMenu   = "settings"
	yesNotify      = "yesNotify"
	noNotify       = "noNotify"
	yesStrategy    = "yesStrategy"
	noStrategy     = "noStrategy"
	offBot         = "stop"
	cancelComm     = "cancel"
)

// Cancel command handler: returns a previous step
func (b *Bot) CancelHandler(ChatId int64) {
	// Floy Warshall branch
	if strings.Contains(UserHistory[ChatId], "FW") {
		if UserHistory[ChatId][3:] == "0" {
			UserHistory[ChatId] = "strategies"
			kb := b.StrategiesKb()
			b.EditAndSend(&kb, strategiesMessage, ChatId)
			return
		} else if UserHistory[ChatId][3:] == "1" || UserHistory[ChatId][3:] == "-1" {
			UserHistory[ChatId] = "FW_0"
			kb := b.YesNoStrategyKb()
			txt := fmt.Sprintf(strategyPower, "Floyd Warshall")
			b.EditAndSend(&kb, txt, ChatId)
			return
		}
		// Strategies
	} else if strings.Contains(UserHistory[ChatId], "strategies") {
		fmt.Println(UserHistory[ChatId])
		kb, txt := b.GetMenuMessage(ChatId)
		b.EditAndSend(&kb, txt, ChatId)
		return
	}
}

// Get main menu keyboard
func (b *Bot) GetMenuMessage(ChatId int64) (tgbotapi.InlineKeyboardMarkup, string) {
	UserHistory[ChatId] = ""
	kb := b.MainKb()
	txt := fmt.Sprintf(menuMessage, 123.1, 0.01)
	return kb, txt
}

// Edit last message and send with current new txt and kb.
func (b *Bot) EditAndSend(kb *tgbotapi.InlineKeyboardMarkup, txt string, ChatId int64) {
	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.Dlg[ChatId].ChatId,
			MessageID:   b.Dlg[ChatId].MessageId,
			ReplyMarkup: kb,
		},
		DisableWebPagePreview: true,
		Text:                  txt,
		ParseMode:             "markdown",
	}
	b.Bot.Send(msg)
}

// Print error in console and send Error text.
func (b *Bot) PrintAndSendError(err error, ChatId int64) {
	fmt.Println(err)
	b.SendMessage(err.Error(), ChatId, nil)
}

// Send Message from tgbotapi with markdown style and current kb.
func (b Bot) SendMessage(txt string, ChatId int64, kb interface{}) {
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	msg.DisableWebPagePreview = true
	b.Bot.Send(msg)
	b.Dlg[ChatId].MessageId++
}

// Send Message from tgbotapi with markdown style and current kb.
func (b Bot) SendNotifyMessage(txt string, ChatId int64, kb interface{}) {
	if ok, n := b.Members[ChatId]; !ok {
		log.Println("Haven't got this user!")
	} else if n {
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = kb
		msg.DisableWebPagePreview = true
		b.Bot.Send(msg)
		b.Dlg[ChatId].MessageId++
	}
}

// Write user's chosen to members.json(gitingore).
func (b *Bot) WriteToJson(ChatId int64, flag bool) {
	data := models.Members{}
	var i = 0
	for k, v := range b.Members {
		data.M = append(data.M, &models.User{
			ChatId:       k,
			Notification: v,
		})
		i++
	}

	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("members.json", file, os.ModePerm)
}
