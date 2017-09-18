package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/proshik/githubstatbot/github"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	client *github.Client
	*github.OAuth
}

func NewBot(token string, debug bool, ghClient *github.Client, clientId string, clientSecret string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	log.Printf("Authorized for account %s", bot.Self.UserName)

	return &Bot{bot, ghClient, &github.OAuth{ClientId: clientId, ClientSecret: clientSecret}}, nil
}
