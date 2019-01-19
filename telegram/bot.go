package telegram

import (
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

var (
	BotName            string
	RedirectBotAddress = "https://t.me/"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	oAuth      *github.OAuth
	tokenStore storage.AccessTokenStorage
	stateStore *storage.StateStore
}

func NewBot(
	token string,
	debug bool,
	tokenStore storage.AccessTokenStorage,
	stateStore *storage.StateStore,
	oAuth *github.OAuth) (*Bot, error) {

	// authorize telegram bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	// set debug mode for bot
	bot.Debug = debug
	// fill botName and Telegram bot URL
	BotName = bot.Self.UserName
	RedirectBotAddress += BotName

	log.Printf("Authorized telegram bot for account %s", bot.Self.UserName)
	return &Bot{bot, oAuth, tokenStore, stateStore}, nil
}
