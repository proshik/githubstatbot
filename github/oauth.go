package github

import (
	"fmt"
)

type OAuth struct {
	ClientId     string
	ClientSecret string
}

func NewOAuth(clientId string, clientSecret string) *OAuth {
	return &OAuth{
		clientId,
		clientSecret,
	}
}

func (oAuth *OAuth) BuildAuthUrl(state string) string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&state=%s", oAuth.ClientId, state)
}
