package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/aliyevazam/telegram-bot/bot"
	"github.com/aliyevazam/telegram-bot/config"
	"github.com/aliyevazam/telegram-bot/pkg/db"
	"github.com/aliyevazam/telegram-bot/storage"
)

func main() {
	cfg := config.LoadConfig(".")

	psqlConn, err := db.ConnectToDb(cfg)
	if err != nil {
		log.Printf("failed to connect database: %v", err)
	}

	strg := storage.NewStoragePg(psqlConn)

	botHandler := bot.New(cfg, strg)

	botHandler.Start()
}
