package api

import (
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
)

type Handler struct {
	oAuth      *github.OAuth
	tokenStore AccessTokenStorage
	stateStore *storage.StateStore
	bot        *telegram.Bot
}

type AccessTokenStorage interface {
	Get(chatId int64) string
	Add(chatId int64, accessToken string)
	Delete(key int64)
}

func New(
	OAuth *github.OAuth,
	tokenStore AccessTokenStorage,
	stateStore *storage.StateStore,
	bot *telegram.Bot) Handler {
	return Handler{OAuth, tokenStore, stateStore, bot}
}
