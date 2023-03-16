package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var enterCommand = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Информация о погоде"),
		tgbotapi.NewKeyboardButton("Статистика"),
	),
)
var statsCommand = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Статистика моих запросов"),
		tgbotapi.NewKeyboardButton("Мой первый запрос"),
	),
)
