package api

import (
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
)

type Handler struct {
	oAuth      *github.OAuth
	tokenStore *storage.TokenStore
	stateStore *storage.StateStore
	bot        *telegram.Bot
}

func New(
	OAuth *github.OAuth,
	tokenStore *storage.TokenStore,
	stateStore *storage.StateStore,
	bot *telegram.Bot) Handler {
	return Handler{OAuth, tokenStore, stateStore, bot}
}
