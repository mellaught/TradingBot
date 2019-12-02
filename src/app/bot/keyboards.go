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
func (b *Bot) YesStrategyKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ON ✔️", yesStrategy),
		),
	)
}

// YesNoTradingKb: ON or OFF Notifications/BOT.
func (b *Bot) NoStrategyKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("OFF ❌", noStrategy),
		),
	)
}

// MainKb: Notifications ON/OFF, Trading Bot Stop/RUN.
func (b *Bot) MainKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Notifications 🔔", notifyCommand),
			tgbotapi.NewInlineKeyboardButtonData("Trading 📊", strategyCommand),
		),
	)
}

// StrategiesKb: choose strategy
func (b *Bot) StrategiesKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Floyd Warshall 🖍️", fwCommand),
			tgbotapi.NewInlineKeyboardButtonData("Triangular 🖍️", trigCommand),
		),
	)
}

// StrategiesKb: choose exhcange
func (b *Bot) ExchangesKb() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Binance 💱", binanceCommand),
			tgbotapi.NewInlineKeyboardButtonData("Poloniex 💱", poloniexCommand),
		),
	)
}

// CancelKeyboard for cancel step.
func (b *Bot) CancelKeyboard(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cancel", cancelCommand),
		),
	)
}
