package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestReposFoundByUser(t *testing.T) {
	startServer()
	defer teardown()

	var user = "proshik"
	mux.HandleFunc(fmt.Sprintf("/users/%s/repos", user), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"repo1", "fork":false}, {"name":"repo2","fork":false}, {"name":"repo3","fork":true}]`)
	})

	client := &Client{client: client}
	repos, err := client.Repos(user)
	if err != nil {
		t.Errorf("Repos return erroir: %v", err)
	}

	want := []*Repo{{Name: String("repo1")}, {Name: String("repo2")}}
	if !reflect.DeepEqual(want, repos) {
		t.Errorf("Repos want: %+v, returned: %+v", want, repos)
	}
}

func TestReposNotFoundByUser(t *testing.T) {
	startServer()
	defer teardown()

	var user = "proshik"
	mux.HandleFunc(fmt.Sprintf("/users/%s/repos", user), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found 404")
	})

	client := &Client{client: client}
	repos, err := client.Repos(user)
	if err == nil {
		t.Errorf("Repos must be return error, but found result: %v", repos)
	}
}

/*
// call read service with token=GITHUB_TOKEN from environment variable
func TestReposSuccess(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	user := "proshik";
	github, err := NewClient(token)

	_, err = github.Repos(user)
	if err != nil {
		t.Errorf("Not found not one repository by username=%s", user)
	}
}
*/
