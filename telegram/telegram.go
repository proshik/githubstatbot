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

	start := make(chan tgbotapi.Update)
	language := make(chan tgbotapi.Update)
	bot_res := make(chan tgbotapi.Chattable)

	go func() {
		for {
			select {
			case update := <-start:
				bot_res <- startCommand(&update)
			case update := <-language:
				bot_res <- languageCommand(&update, b.GitHubClient)
			}
		}
	}()

	go func() {
		for res := range bot_res {
			fmt.Println("handle bot_res")
			b.Bot.Send(res)
		}
	}()

	for update := range updates {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				start <- update
				continue
			case "language":
				language <- update
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unrecognized command. Say what?")
				b.Bot.Send(msg)
			}
		} else {
			start <- update
		}
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
