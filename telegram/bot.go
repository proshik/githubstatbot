package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
)

type Bot struct {
	bot        *tgbotapi.BotAPI
	oAuth      *github.OAuth
	tokenStore *storage.TokenStore
	stateStore *storage.StateStore
}

//type AccessToken interface {
//	Get(chatId int64) (string, error)
//	Add(chatId int64, accessToken string)
//}

func NewBot(
	token string,
	debug bool,
	tokenStore *storage.TokenStore,
	stateStore *storage.StateStore,
	oAuth *github.OAuth) (*Bot, error) {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	log.Printf("Authorized telegram bot for account %s", bot.Self.UserName)

	return &Bot{bot, oAuth, tokenStore, stateStore}, nil
}
