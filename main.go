package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/api"
	"github.com/proshik/githubstatbot/config"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	// config logging
	log.SetOutput(os.Stdout)
	// init connect to db(boltDB)
	//db := storage.New(cfg.DbPath)
	db := storage.NewPostgres(cfg.DbUrl)
	// create storage for generated statuses for request to github.com
	stateStore := storage.NewStateStore()
	// create oAuth object
	oAuth := github.NewOAuth(cfg.GitHubClientId, cfg.GitHubClientSecret)
	// create Telegram Bot object
	bot, err := telegram.NewBot(cfg.TelegramToken, false, db, stateStore, oAuth)
	if err != nil {
		log.Panic(err)
	}
	// run major method for read updates messages from telegram
	go bot.ReadUpdates()
	// initialize handler
	basicAuth := &api.BasicAuth{Username: cfg.AuthBasicUsername, Password: cfg.AuthBasicPassword}
	handler := api.New(oAuth, db, stateStore, bot, basicAuth, cfg.StaticFilesDir)
	// configuration router
	router := httprouter.New()
	router.GET("/", handler.Index)
	router.GET("/version", handler.Version)
	router.GET("/github_redirect", handler.GitHubRedirect)

	//Run HTTPS server
	http.ListenAndServe(":"+cfg.Port, router)
}
