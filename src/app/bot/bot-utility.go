package bot

import (
	"TradingBot/src/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand     = "start"
	getMainMenu      = "home"
	exchangesCommand = "exchanges"
	strategyCommand  = "strategy"
	binanceCommand   = "binance"
	poloniexCommand  = "poloniex"
	trigCommand      = "tri"
	fwCommand        = "fw"
	notifyCommand    = "notify"
	settingsMenu     = "settings"
	yesNotify        = "yesNotify"
	noNotify         = "noNotify"
	yesStrategy      = "yesStrategy"
	noStrategy       = "noStrategy"
	offBot           = "stop"
	cancelCommand    = "cancel"
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
			// IF User turned ON strategy
			if _, ok := b.MembersStrategy[ChatId][b.UserStrategy[ChatId]]["FW"]; ok {
				kb := b.YesStrategyKb()
				txt := fmt.Sprintf(strategyPower, "Floyd Warshall")
				b.EditAndSend(&kb, txt, ChatId)
			} else {
				kb := b.NoStrategyKb()
				txt := fmt.Sprintf(strategyPower, "Floyd Warshall")
				b.EditAndSend(&kb, txt, ChatId)
			}

			return
		}

		// Strategies
	} else if strings.Contains(UserHistory[ChatId], "exchanges") {
		fmt.Println(UserHistory[ChatId])
		kb, txt := b.GetMenuMessage(ChatId)
		b.EditAndSend(&kb, txt, ChatId)
		return
	} else if strings.Contains(UserHistory[ChatId], "strategies") {
		fmt.Println(UserHistory[ChatId])
		kb, txt := b.GetMenuMessage(ChatId)
		b.EditAndSend(&kb, txt, ChatId)
		return
	}
}

// Strategy handler: turn on or off strategy
func (b *Bot) StrategyHandler(exchange, strategy string, on bool, ChatId int64) {
	// Turn ON strategy for User
	if on {
		if e, ok := b.MembersStrategy[ChatId]; ok {
			if s, ok := e[exchange]; ok {
				fmt.Println(s)
				log.Printf("Strategy %s on %s exchang for User %d is already ON", strategy, exchange, ChatId)
				return
			}
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			b.MembersStrategy[ChatId][exchange][strategy] = &Strategy{
				Ctx:    &ctx,
				Cancel: cancel,
			}
			// Send notification to StrategyHandler in App
			notify := ExchangeStrategy{
				ChatId:   ChatId,
				Strategy: strategy,
				Exchange: exchange,
				Ctx:      &ctx,
			}

			b.RunStrategy <- notify
		}
		// Turn OFF strategy for User
	} else {
		if e, ok := b.MembersStrategy[ChatId][exchange][strategy]; ok {
			e.Cancel()
			log.Printf("Strategy %s OFF on %s exchang for User %d", strategy, exchange, ChatId)
		} else {
			log.Printf("Strategy %s on %s exchang for User %d doesn't ON", strategy, exchange, ChatId)
			return
		}
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
