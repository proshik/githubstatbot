package github

import (
	"context"
	gh "github.com/google/go-github/github"
)

type Repo struct {
	Name            *string `json:"name"`
	Language        *string `json:"language"`
	StargazersCount *int    `json:"stargazers_count"`
	ForksCount      *int    `json:"forks_count"`
}

func (github *Client) Repos(user string) ([]*Repo, error) {
	ctx := context.Background()
	opt := gh.RepositoryListOptions{Sort: "updated", ListOptions: gh.ListOptions{PerPage: 100}}

	repos, _, err := github.client.Repositories.List(ctx, user, &opt)
	if err != nil {
		return nil, err
	}

	result := make([]*Repo, 0)
	for _, r := range repos {
		//skip forked repositories
		if *r.Fork {
			continue
		}
		result = append(result, &Repo{r.Name, r.Language, r.StargazersCount, r.ForksCount})
	}

	return result, nil
}

func (github *Client) Repo(user string, repoName string) (*Repo, error) {
	ctx := context.Background()

	r, _, err := github.client.Repositories.Get(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return &Repo{r.Name, r.Language, r.StargazersCount, r.ForksCount}, nil
}
