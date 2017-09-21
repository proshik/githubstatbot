package main

import (
	"os"
	"log"
	"net/http"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/telegram"
	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/api"
	"github.com/proshik/githubstatbot/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Panic("Port is empty")
	}

	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		log.Panic("ClientId or clientSecret is empty")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Panic("Telegram token is empty")
	}

	tokenStore := storage.NewTokenStore()
	stateStore := storage.NewStateStore()

	oAuth := github.NewOAuth(clientId, clientSecret)

	bot, err := telegram.NewBot(telegramToken, false, tokenStore, stateStore, oAuth)
	if err != nil {
		log.Panic(err)
	}

	handler := api.New(oAuth, tokenStore, stateStore, bot)

	go bot.ReadUpdates()

	router := httprouter.New()
	router.GET("/github_redirect", handler.GitHubRedirect)

	log.Println("Service is waiting for requests...")

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
