package github

import (
	"github.com/google/go-github/github"
	"context"
)

func (github *Client) CommitActivity(user string, repoName string) ([]*github.WeeklyCommitActivity, error) {
	ctx := context.Background()

	activity, _, err := github.client.Repositories.ListCommitActivity(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return activity, nil
}
