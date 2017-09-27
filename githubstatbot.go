package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/api"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"log"
	"net/http"
	"os"
)

//For run:
//env PORT=8080 DB_PATH=/data/githubstatbot/boltdb.db GITHUB_CLIENT_ID= GITHUB_CLIENT_SECRET= TELEGRAM_TOKEN= go run githubstatbot.go
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Panic("Port is empty")
	}

	path := os.Getenv("DB_PATH")
	if path == "" {
		log.Panic("DB path is empty")
	}

	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Panic("ClientId or clientSecret is empty")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Panic("Telegram token is empty")
	}

	db := storage.New(path)
	stateStore := storage.NewStateStore()
	oAuth := github.NewOAuth(clientID, clientSecret)

	bot, err := telegram.NewBot(telegramToken, false, db, stateStore, oAuth)
	if err != nil {
		log.Panic(err)
	}
	go bot.ReadUpdates()

	handler := api.New(oAuth, db, stateStore, bot)
	router := httprouter.New()
	router.GET("/", handler.Index)
	router.GET("/github_redirect", handler.GitHubRedirect)

	log.Println("Service is waiting for requests...")

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
