package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TokenStoreMock struct {
	store map[int64]string
}

func (s *TokenStoreMock) Get(chatId int64) (string, error) {
	return s.store[chatId], nil
}

func (s *TokenStoreMock) Add(chatId int64, accessToken string) error {
	s.store[chatId] = accessToken
	return nil
}

func (s *TokenStoreMock) Delete(key int64) error {
	delete(s.store, key)
	return nil
}

var (
	oAuthMock      = github.NewOAuth("clientId", "clientSecret")
	tokenStoreMock = &TokenStoreMock{}
	stateStoreMock = storage.NewStateStore()
	bot, _         = telegram.NewBot("telegramToken", false, tokenStoreMock, stateStoreMock, oAuthMock)
	basicAuth      = &BasicAuth{"username", "password"}

	h = New(
		oAuthMock,
		tokenStoreMock,
		stateStoreMock,
		bot,
		basicAuth,
	)
)

func TestIndex(t *testing.T) {
	router := httprouter.New()
	router.GET("/", h.Index)

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	expectedText := "<html><body>Welcome to GitHubStat Bot!</body></html>"
	actualText := string(greeting)
	if expectedText != string(greeting) {
		t.Fatalf(
			"Wrong text on Index page '%s', expected '%s'",
			actualText, expectedText,
		)
	}
}
