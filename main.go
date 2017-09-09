package main

import (
	"os"
	"log"
	"github.com/proshik/githublangbot/client"
	"github.com/proshik/githublangbot/telegram"
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

	github, err := client.NewGithub(gitHubToken)
	if err != nil {
		log.Panic(err)
	}

	bot, err := telegram.NewBot(telegramToken, false, github)
	if err != nil {
		log.Panic(err)
	}

	bot.ReadUpdates()
}
