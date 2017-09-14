package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/proshik/githublangbot/github"
)

type Bot struct {
	Bot          *tgbotapi.BotAPI
	Client *github.Client
}

func NewBot(token string, debug bool, ghClient *github.Client) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	log.Printf("Authorized for account %s", bot.Self.UserName)

	return &Bot{bot, ghClient}, nil
}
