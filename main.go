package main

import (
	"os"
	"log"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/telegram"
)

func main() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Panic("Telegram token is empty")
	}

	gitHubToken := os.Getenv("GITHUB_TOKEN")
	if gitHubToken == "" {
		log.Panic("GitHub token is empty")
	}

	client, err := github.NewClient(gitHubToken)
	if err != nil {
		log.Panic(err)
	}

	bot, err := telegram.NewBot(telegramToken, false, client)
	if err != nil {
		log.Panic(err)
	}

	bot.ReadUpdates()
}
