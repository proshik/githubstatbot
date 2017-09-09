package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"bytes"
	"fmt"
	"github.com/proshik/githublangbot/client"
)

func (b *Bot) ReadUpdates() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.Bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		var msg tgbotapi.MessageConfig

		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "start":
				msg = startCommand(&update)
			case "language":
				msg = languageCommand(&update, b.GitHubClient)
			default:
				msg = sendCommandInfo(&update)
			}
		} else {
			msg = sendCommandInfo(&update)
		}

		b.Bot.Send(msg)
	}
}
func startCommand(update *tgbotapi.Update) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, createStartText())
	msg.ParseMode = "markdown"

	return msg
}

func languageCommand(update *tgbotapi.Update, github *client.GitHub) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	if update.Message.CommandArguments() != "" {
		r, err := github.Languages(update.Message.CommandArguments())
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error on request")
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, createLangStatText(r))
		msg.ParseMode = "markdown"
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Argument is empty")
	}

	return msg
}

func sendCommandInfo(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "default message")
}


func createStartText() string {
	buf := bytes.NewBufferString("Telegram bot written in GO. This bot show GitHub languages info by account\n")

	buf.WriteString("\n")
	buf.WriteString("You can control me by sending these commands:\n\n")
	buf.WriteString("*/languages <github_account_name>* - list languages for the all repositories\n")

	return buf.String()
}
func createLangStatText(statistics []*client.Languages) string {
	buf := bytes.NewBufferString("")

	for _, l := range statistics {
		buf.WriteString(fmt.Sprintf("*%s* %.1f%%\n", l.Language, l.Percentage))
	}

	return buf.String()
}
