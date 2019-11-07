package bot

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand = "start"
	getMainMenu  = "home"
	subscribe    = "subscribe"
	unsubscribe  = "unsubscribe"
	settingsMenu = "settings"
	cancelComm   = "cancel"
	yescommand   = "yes"
	nocommand    = "not"
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
