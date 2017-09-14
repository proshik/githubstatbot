package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"bytes"
	"fmt"
	_"github.com/proshik/githublangbot/github"
	"sync"
	"sort"
	gh "github.com/google/go-github/github"
)

type Languages struct {
	Language   string
	Percentage float32
}

func (b *Bot) ReadUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.Bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	start := make(chan tgbotapi.Update)
	info := make(chan tgbotapi.Update)
	language := make(chan tgbotapi.Update)
	bot_res := make(chan tgbotapi.Chattable)

	go func() {
		for {
			select {
			case update := <-start:
				bot_res <- startCommand(&update)
			case update := <-info:
				bot_res <- sendCommandInfo(&update)
			case update := <-language:
				done := languageCommand(&update, b)
				bot_res <- done
			}
		}
	}()

	go func() {
		for res := range bot_res {
			b.Bot.Send(res)
		}
	}()

	for update := range updates {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				start <- update
			case "language":
				language <- update
			default:
				info <- update
			}
		} else {
			start <- update
		}
	}
}

func startCommand(update *tgbotapi.Update) tgbotapi.MessageConfig {
	buf := bytes.NewBufferString("Telegram bot written in GO. This bot show GitHub languages info by account\n")

	buf.WriteString("\n")
	buf.WriteString("You can control me by sending these commands:\n\n")
	buf.WriteString("*/languages <user>* - list languages for the all repositories\n")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, buf.String())
	msg.ParseMode = "markdown"

	return msg
}

func sendCommandInfo(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "default message")
}

func languageCommand(update *tgbotapi.Update, github Repository) tgbotapi.MessageConfig {
	user := update.Message.CommandArguments()

	repos, err := github.Repos(user)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Not found info by user with name="+user)
	}

	wg := sync.WaitGroup{}
	languageChan := make(chan map[string]int)
	for _, repo := range repos {
		wg.Add(1)
		go func(wg *sync.WaitGroup, repo *gh.Repository) {
			defer wg.Done()

			lang, err := github.Language(user, *repo.Name)
			if err != nil {
				log.Printf("Error on request language for user=%s, repo=%s", user, *repo.Name)
			}
			languageChan <- lang
		}(&wg, repo)
	}

	go func() {
		wg.Wait()
		close(languageChan)
	}()

	statistics := make(map[string]int)
	for stat := range languageChan {
		for k, v := range stat {
			statistics[k] = statistics[k] + v
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, createLangStatText(calcPercentages(statistics)))
	msg.ParseMode = "markdown"

	return msg
}

func calcPercentages(languages map[string]int) []*Languages {
	result := make([]*Languages, 0)

	var totalSum float32
	for _, v := range languages {
		totalSum += float32(v)
	}

	for key, value := range languages {
		percent := float32(value) * (float32(100) / totalSum)
		result = append(result, &Languages{key, round(percent, 0.1)})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Percentage > result[j].Percentage })

	return result
}

func round(x, unit float32) float32 {
	if x > 0 {
		return float32(int32(x/unit+0.5)) * unit
	}
	return float32(int32(x/unit-0.5)) * unit
}

func createLangStatText(statistics []*Languages) string {
	buf := bytes.NewBufferString("")

	for _, l := range statistics {
		buf.WriteString(fmt.Sprintf("*%s* %.1f%%\n", l.Language, l.Percentage))
	}

	return buf.String()
}
