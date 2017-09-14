package github

import "context"

func (github *Client) Language(user string, repoName string) (map[string]int, error) {
	ctx := context.Background()

	lang, _, err := github.client.Repositories.ListLanguages(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return lang, nil
}
