package bot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// YesNoKB: ON or OFF Notifications/BOT.
func (b *Bot) YesNoKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ON✔️", yescommand),
			tgbotapi.NewInlineKeyboardButtonData("OFF❌", nocommand),
		),
	)
}

// MainKb: Notifications ON/OFF, Trading Bot Stop/RUN.
func (b *Bot) MainKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Notifications🔔", notifyCommand),
			tgbotapi.NewInlineKeyboardButtonData("Trading📊", tradingCommand),
		),
	)
}
