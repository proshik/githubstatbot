package client

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"context"
)

type GitHub struct {
	Client *github.Client
}

func NewGitHub(token string) (*GitHub, error) {
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

func (github *GitHub) Repo(user string, repoName string) (*github.Repository, error) {
	ctx := context.Background()

	repo, _, err := github.Client.Repositories.Get(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (github *GitHub) Language(user string, repoName string) (map[string]int, error) {
	ctx := context.Background()

	lang, _, err := github.Client.Repositories.ListLanguages(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return lang, nil
}

func (github *GitHub) CommitActivity(user string, repoName string) ([]*github.WeeklyCommitActivity, error) {
	ctx := context.Background()

	activity, _, err := github.Client.Repositories.ListCommitActivity(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
