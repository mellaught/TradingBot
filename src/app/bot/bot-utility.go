package bot

import (
	"TradingBot/src/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand   = "start"
	getMainMenu    = "home"
	tradingCommand = "trading"
	notifyCommand  = "notify"
	settingsMenu   = "settings"
	yesNotify      = "yes"
	noNotify       = "no"
	offBot         = "stop"
)

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
