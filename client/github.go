package client

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"context"
	_"log"
	_"sort"
)

type GitHub struct {
	Client *github.Client
}

func NewGithub(token string) (*GitHub, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)

	return &GitHub{c}, nil
}

func (github *GitHub) Repos(user string) ([]*github.Repository, error) {
	ctx := context.Background()

	repos, _, err := github.Client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (github *GitHub) Language(user string, repoName string) (map[string]int, error) {
	ctx := context.Background()

	lang, _, err := github.Client.Repositories.ListLanguages(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return lang, nil
}
