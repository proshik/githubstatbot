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

	page := 1
	result := make([]*Repo, 0)
	for {
		opt := gh.RepositoryListOptions{Sort: "updated", ListOptions: gh.ListOptions{PerPage: 100, Page: page}}
		list, resp, err := github.client.Repositories.List(ctx, user, &opt)
		if err != nil {
			return nil, err
		}
		// fill result slice
		for _, r := range list {
			if *r.Fork {
				continue
			}
			result = append(result, &Repo{r.Name, r.Language, r.StargazersCount, r.ForksCount})
		}
		// check on exist next page
		if page <= resp.LastPage {
			page++
		} else {
			break
		}
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
