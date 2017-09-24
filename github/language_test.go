package github

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLanguagesFound(t *testing.T) {
	startServer()
	defer teardown()

	user := "proshik"
	repoNames := []string{"repo1", "repo2"}
	for _, name := range repoNames {
		mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/languages", user, name), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"GO":1314}, {"Java":32342}`)
		})
	}

	client := &Client{client: client}
	for _, name := range repoNames {
		lang, err := client.Language(user, name)
		if err != nil {
			t.Errorf("Languages return error: %v", err)
		}
		if len(lang) == 0 {
			t.Errorf("Lang from response is empty: %+v", lang)
		}
	}
}
