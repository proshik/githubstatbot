package api

import (
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"net/http"
)

type Handler struct {
	oAuth      *github.OAuth
	tokenStore storage.AccessTokenStorage
	stateStore *storage.StateStore
	bot        *telegram.Bot
}

func New(
	OAuth *github.OAuth,
	tokenStore storage.AccessTokenStorage,
	stateStore *storage.StateStore,
	bot *telegram.Bot) Handler {
	return Handler{OAuth, tokenStore, stateStore, bot}
}

func (h *Handler) RedirectToHttps(w http.ResponseWriter, r *http.Request) {
	newURI := "https://" + r.Host + r.URL.String()
	http.Redirect(w, r, newURI, http.StatusFound)
}
