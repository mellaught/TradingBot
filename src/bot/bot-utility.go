package bot

import (
	"TradingBot/src/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand   = "start"
	getMainMenu    = "home"
	tradingCommand = "trading"
	notifyCommand  = "notify"
	settingsMenu   = "settings"
	yescommand     = "yes"
	nocommand      = "not"
)

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

func (b *Bot) PrintAndSendError(err error, ChatId int64) {
	fmt.Println(err)
	b.SendMessage("Error", ChatId, nil)
}

func (b Bot) SendMessage(txt string, ChatId int64, kb interface{}) {
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	msg.DisableWebPagePreview = true
	b.Bot.Send(msg)
	b.Dlg[ChatId].MessageId++
}

// Write user's chosen to members.json(gitingore).
func (b *Bot) WriteToJson(ChatId int64, flag bool) {
	data := &models.Members{}
	for i, m := range b.Members.M {
		data.M[i] = &models.User{
			ChatId:       m.ChatId,
			Notification: m.Notification,
		}
	}

	file, _ := json.Marshal(data)
	_ = ioutil.WriteFile("members.json", file, os.ModePerm)
}
