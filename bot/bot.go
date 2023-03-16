package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aliyevazam/telegram-bot/config"
	"github.com/aliyevazam/telegram-bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	cfg     config.Config
	storage storage.StorageI
	bot     *tgbotapi.BotAPI
}

func New(cfg config.Config, strg storage.StorageI) BotHandler {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		fmt.Println("error", err)
		log.Panic(err)
	}
	bot.Debug = true
	return BotHandler{
		cfg:     cfg,
		storage: strg,
		bot:     bot,
	}
}

func (h *BotHandler) Start() {
	log.Printf("Authorized on account %s", h.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		go h.HandleBot(update)
	}
}

func (h *BotHandler) HandleBot(update tgbotapi.Update) {
	if update.Message.Command() == "start" {
		err := h.HandleEnterCommand(update)
		if err != nil {
			log.Printf("Error while handler enter command: %v", err)
			return
		}
	} else if update.Message.Text != "" {
		if update.Message.Text == "Информация о погоде" {
			err := h.HandleEnterCountry(update)
			if err != nil {
				log.Printf("Error while send to handle enter country: %v", err)
				return
			}
		} else if update.Message.Text == "Статистика" {
			err := h.HandleEnterStats(update)
			if err != nil {
				log.Printf("Error while get Stats: %v", err)
				return
			}
		} else if update.Message.Text == "Статистика моих запросов" {
			response, err := h.storage.GetRequest(update.Message.From.ID)
			if err != nil {
				log.Printf("Error while get requests: %v", err)
				return
			}
			h.SendDataToUserArray(update, response)

		} else if update.Message.Text == "Мой первый запрос" {
			response, err := h.storage.GetFirstRequest(update.Message.From.ID)
			if err != nil {
				log.Printf("Error while get first request: %v", err)
				return
			}
			h.SendDataToUser(update, response)
		} else {
			response, err := getInfo(update.Message.Text, h.cfg.ApiToken)
			if err != nil {
				log.Printf("Error while get info: %v", err)
				return
			}
			if response == "Ошибка в имени города" {
				h.SendMessage(update, "Ошибка в имени города пожалуйста проверте")
				return
			}
			err = h.storage.CreateRequest(update.Message.Text, update.Message.From.ID)
			if err != nil {
				log.Printf("Error while create request: %v", err)
				return
			}
			h.SendMessage(update, response)
		}
	}

}

func getInfo(query string, apiKey string) (string, error) {
	// Обработка запроса для поиска информации
	// В этом примере мы будем искать погоду в указанном городе
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", query, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while get information from weather: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Чтение ответа и обработка данных
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var weather struct {
		Name string `json:"name"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}
	err = json.Unmarshal(data, &weather)
	if err != nil {
		return "", err
	}
	if weather.Main.Temp == 0.0 {
		return "Ошибка в имени города", nil
	}

	message := fmt.Sprintf("Температура в %s: %.1f°C", weather.Name, weather.Main.Temp)
	return message, nil
}
