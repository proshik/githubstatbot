package github

import "context"

func (github *Client) User() (string, error) {
	ctx := context.Background()

	user, _, err := github.client.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}

	return *user.Login, nil
}
