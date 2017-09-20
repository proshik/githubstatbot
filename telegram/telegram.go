package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"bytes"
	"fmt"
	"sync"
	"sort"
	"github.com/proshik/githubstatbot/github"
	"math/rand"
	"strconv"
)

type Language struct {
	Title      string
	Percentage float32
}

func (b *Bot) ReadUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	start := make(chan tgbotapi.Update)
	info := make(chan tgbotapi.Update)
	auth := make(chan tgbotapi.Update)
	language := make(chan tgbotapi.Update)
	bot_res := make(chan tgbotapi.Chattable)

	go func() {
		for {
			select {
			case u := <-start:
				bot_res <- startCommand(&u)
			case u := <-info:
				bot_res <- infoCommand(&u)
			case u := <-auth:
				bot_res <- authCommand(&u, b)
			case u := <-language:
				done := languageCommand(&u, b)
				bot_res <- done
			}
		}
	}()

	go func() {
		for res := range bot_res {
			b.bot.Send(res)
		}
	}()

	for update := range updates {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				start <- update
			case "language":
				language <- update
			case "auth":
				auth <- update
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

func infoCommand(update *tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "You must write command")
}

func authCommand(update *tgbotapi.Update, bot *Bot) tgbotapi.Chattable {
	authUrl := bot.oAuth.BuildAuthUrl(strconv.Itoa(int(update.Message.Chat.ID)))

	return tgbotapi.NewMessage(update.Message.Chat.ID, authUrl)
}

func languageCommand(update *tgbotapi.Update, bot *Bot) tgbotapi.MessageConfig {
	user := update.Message.CommandArguments()

	token, err := bot.tokenStore.Get(update.Message.Chat.ID)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "You must authorization with github.com, command /auth")
	}

	client, err := github.NewClient(token)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Error on authentication. Please do re-auth with command /auth")
	}

	if user == "" {
		user, err = client.User()
		if err != nil {
			return tgbotapi.NewMessage(update.Message.Chat.ID, "Not found by token. Please do re-auth=")
		}
	}

	repos, err := client.Repos(user)
	if err != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Not found repos for user="+user)
	}

	wg := sync.WaitGroup{}
	languageChan := make(chan map[string]int)
	for _, repo := range repos {
		wg.Add(1)
		go func(wg *sync.WaitGroup, r *github.Repo) {
			defer wg.Done()

			lang, err := client.Language(user, *r.Name)
			if err != nil {
				log.Printf("Error on request language for user=%s, repo=%s", user, *r.Name)
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

func calcPercentages(languages map[string]int) []*Language {
	result := make([]*Language, 0)
	//calculate total sum byte by all languages
	var totalSum float32
	for _, v := range languages {
		totalSum += float32(v)
	}

	var totalByteOtherLanguages int
	for key, value := range languages {
		repoPercent := float32(value) * (float32(100) / totalSum)
		roundRepoPercent := round(repoPercent, 0.1)
		if roundRepoPercent >= 0.1 {
			result = append(result, &Language{key, roundRepoPercent})
		} else {
			totalByteOtherLanguages += value
		}
	}
	//sort found languages by percentage
	sort.Slice(result, func(i, j int) bool { return result[i].Percentage > result[j].Percentage })
	//calculate percentage for language with less then 0.1% from total count
	if totalByteOtherLanguages != 0 {
		percent := round(float32(totalByteOtherLanguages)*(float32(100)/totalSum), 0.1)
		result = append(result, &Language{"Other languages", percent})
	}

	return result
}

func round(x, unit float32) float32 {
	if x > 0 {
		return float32(int32(x/unit+0.5)) * unit
	}
	return float32(int32(x/unit-0.5)) * unit
}

func createLangStatText(statistics []*Language) string {
	buf := bytes.NewBufferString("")

	for _, l := range statistics {
		buf.WriteString(fmt.Sprintf("*%s* %.1f%%\n", l.Title, l.Percentage))
	}

	return buf.String()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


