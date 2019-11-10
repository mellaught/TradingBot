package bot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// NotificationsKb: ON or OFF Notifications from TradingBot.
func (b *Bot) NotificationsKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ON", yescommand),
			tgbotapi.NewInlineKeyboardButtonData("OFF", nocommand),
		),
	)
}
