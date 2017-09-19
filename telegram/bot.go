package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/proshik/githubstatbot/github"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	Storage Storage
	client  *github.Client
	*github.OAuth
}

func NewBot(token string, debug bool, clientId string, clientSecret string, ghClient *github.Client) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	storage := NewStorage()

	bot.Debug = debug

	log.Printf("Authorized for account %s", bot.Self.UserName)

	return &Bot{bot, storage, ghClient, &github.OAuth{ClientId: clientId, ClientSecret: clientSecret}}, nil
}
