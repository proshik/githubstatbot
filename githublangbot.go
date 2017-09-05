package main

import (
	"log"
	"os"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/proshik/githublangbot/github"
	"fmt"
	"bytes"
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

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized for account %s", bot.Self.UserName)

	gh := &github.GHAuth{
		Token: gitHubToken,
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		var msg tgbotapi.MessageConfig

		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "languages":
				if update.Message.CommandArguments() != "" {
					r, err := gh.LangStatistic(update.Message.CommandArguments())
					if err != nil {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error on request")
					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, createLangStatText(r))
					msg.ParseMode = "markdown"
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Argument is empty")
				}
			default:
				msg = sendCommandInfo(update)
			}
		} else {
			msg = sendCommandInfo(update)
		}

		bot.Send(msg)
	}
}
func createLangStatText(statistics []*github.LangStatistic) string {
	buf := bytes.NewBufferString("")

	for _, l := range statistics {
		buf.WriteString(fmt.Sprintf("*%s* %.1f%%\n", l.Language, l.Percentage))
	}

	return buf.String()
}

func sendCommandInfo(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "default message")
}
