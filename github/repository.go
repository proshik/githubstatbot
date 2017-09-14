package github

import (
	"github.com/google/go-github/github"
	"context"
)

type Repo struct {
	Name *string `json:"name, omitempty"`
}

func (github *Client) Repos(user string) ([]*Repo, error) {
	ctx := context.Background()

	repos, _, err := github.client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*Repo, 0)
	for _, r := range repos {
		result = append(result, &Repo{r.Name})
	}

	return result, nil
}

func (github *Client) Repo(user string, repoName string) (*Repo, error) {
	ctx := context.Background()

	repo, _, err := github.client.Repositories.Get(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return &Repo{repo.Name}, nil
}

func (github *Client) CommitActivity(user string, repoName string) ([]*github.WeeklyCommitActivity, error) {
	ctx := context.Background()

	activity, _, err := github.client.Repositories.ListCommitActivity(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
