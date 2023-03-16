package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyevazam/telegram-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) HandleEnterCommand(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Выберите команду")
	msg.ReplyMarkup = enterCommand
	h.storage.GetOrCreate(&storage.User{
		TgID:   update.Message.From.ID,
		TgName: update.Message.From.UserName,
	})
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error while send message from handle enter command: %v", err)
	}
	return nil
}

func (h *BotHandler) HandleEnterCountry(update tgbotapi.Update) error {
	h.SendMessage(update, "Вводите город")
	return nil
}

func (h *BotHandler) HandleEnterStats(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.From.ID, "Статистика")
	msg.ReplyMarkup = statsCommand
	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err)
	}
	return nil
}

func (h *BotHandler) SendMessage(update tgbotapi.Update, message interface{}) {
	fmt.Println("Sendmessageda", message)
	if s, ok := message.(string); ok {
		msg := tgbotapi.NewMessage(update.Message.From.ID, s)
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error while send msg: %v", err)
		}
	}
}

func (h *BotHandler) SendDataToUser(update tgbotapi.Update, message *storage.Request) {
	text := "Ваш первый запрос:\n"
	text += fmt.Sprintf("%s", message)
	symbolToRemove := []string{"&", "{", "}"}
	for _, symbol := range symbolToRemove {
		text = strings.ReplaceAll(text, symbol, "")
	}
	msg := tgbotapi.NewMessage(update.Message.From.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error while send msg: %v", err)
	}
}

func (h *BotHandler) SendDataToUserArray(update tgbotapi.Update, message []*storage.Request) {
	text := "История ваших запросов:\n"
	for _, req := range message {
		text += fmt.Sprintf("\n%s\n%s\n", req.Query, req.Request_time)
	}
	msg := tgbotapi.NewMessage(update.Message.From.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error while send msg: %v", err)
	}
}
