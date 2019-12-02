package bot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// YesNoNotifyKb: ON or OFF Notifications/BOT.
func (b *Bot) YesNoNotifyKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ON ✔️", yesNotify),
			tgbotapi.NewInlineKeyboardButtonData("OFF ❌", noNotify),
		),
	)
}

// YesNoStrategyKb: ON or OFF Strategy.
func (b *Bot) YesNoStrategyKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ON ✔️", yesStrategy),
			tgbotapi.NewInlineKeyboardButtonData("OFF ❌", noStrategy),
		),
	)
}

// YesNoTradingKb: ON or OFF Notifications/BOT.
func (b *Bot) YesNoTradingKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("OFF ❌", offBot),
		),
	)
}

// MainKb: Notifications ON/OFF, Trading Bot Stop/RUN.
func (b *Bot) MainKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Notifications 🔔", notifyCommand),
			tgbotapi.NewInlineKeyboardButtonData("Trading 📊", tradingCommand),
		),
	)
}

// MainKb: Notifications ON/OFF, Trading Bot Stop/RUN.
func (b *Bot) StrategiesKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Floyd Warshall 🖍️", fwcommand),
			tgbotapi.NewInlineKeyboardButtonData("Triangular 🖍️", tricommand),
		),
	)
}

// CancelKeyboard for cancel step.
func (b *Bot) CancelKeyboard(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cancel", cancelComm),
		),
	)
}
