package github

import (
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"context"
)

type OAuth struct {
	ClientId string
	ClientSecret string
}

type Client struct {
	client *github.Client
}

func NewClient(token string) (*Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)

	return &Client{c}, nil
}
