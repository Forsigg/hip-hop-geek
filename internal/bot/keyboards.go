package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mainKeybordButtons = [][]tgbotapi.KeyboardButton{
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(TodayButtonText),
		tgbotapi.NewKeyboardButton(TodayReleasesButtonText),
		tgbotapi.NewKeyboardButton(MonthReleasesButtonText),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(SubscribeButtonText),
		tgbotapi.NewKeyboardButton(UnsubscribeButtonText),
		tgbotapi.NewKeyboardButton(CheckSubscribeButtonText),
	),
}

var adminKeyboardButtons = append(mainKeybordButtons, tgbotapi.NewKeyboardButtonRow(
	tgbotapi.NewKeyboardButton(RefreshReleasesButtonText),
	tgbotapi.NewKeyboardButton(TestButtonText),
))

var (
	mainKeyboard  = tgbotapi.NewReplyKeyboard(mainKeybordButtons...)
	adminKeyboard = tgbotapi.NewReplyKeyboard(adminKeyboardButtons...)
)
