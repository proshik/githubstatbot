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
	basicAuth  *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

func New(
	OAuth *github.OAuth,
	tokenStore storage.AccessTokenStorage,
	stateStore *storage.StateStore,
	bot *telegram.Bot,
	basicAuth *BasicAuth) Handler {

	return Handler{OAuth, tokenStore, stateStore, bot, basicAuth}
}

func (h *Handler) RedirectToHttps(w http.ResponseWriter, r *http.Request) {
	newURI := "https://" + r.Host + r.URL.String()
	http.Redirect(w, r, newURI, http.StatusFound)
}
