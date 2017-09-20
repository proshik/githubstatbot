package api

import (
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
)

type Handler struct {
	oAuth      *github.OAuth
	tokenStore *storage.TokenStore
}

func New(OAuth *github.OAuth, tokenStore *storage.TokenStore) Handler {
	return Handler{OAuth, tokenStore}
}
