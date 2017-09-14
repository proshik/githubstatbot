package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/proshik/githublangbot/github"
	"github.com/google/go-github/github"
)

type Repository interface {
	Repos(user string) ([]*github.Repository, error)
	Repo(user string, repoName string) (*github.Repository, error)
	Language(user string, repoName string) (map[string]int, error)
	CommitActivity(user string, repoName string) ([]*github.WeeklyCommitActivity, error)
}

type Bot struct {
	Bot          *tgbotapi.BotAPI
	Repository
}

func NewBot(token string, debug bool, ghClient *github.GitHub) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	log.Printf("Authorized for account %s", bot.Self.UserName)

	return &Bot{bot, ghClient}, nil
}
