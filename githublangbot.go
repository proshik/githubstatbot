package main

import (
	"log"
	"os"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Panic("Token is empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized for account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		bot.Send(msg)
	}
}
