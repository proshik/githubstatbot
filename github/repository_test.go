package github

import (
	"testing"
	"os"
	"github.com/google/go-github/github"
	"net/http"
	"net/http/httptest"
	"net/url"
	"fmt"
	"reflect"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux
	// client is the GitHub client being tested.
	client *github.Client
	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func startServer() *http.ServeMux {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	// github client configured to use test server

	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	client.UploadURL = url

	return mux
}

func teardown() {
	server.Close()
}

func TestReposFoundByUser(t *testing.T) {
	startServer()
	defer teardown()

	var user = "proshik"
	mux.HandleFunc(fmt.Sprintf("/users/%s/repos", user), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"repo1"}, {"name":"repo3"}]`)
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

func TestReposNotFoundByUser(t *testing.T){
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

func TestNewGithubClient(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")

	_, err := NewGitHub(token)
	if err != nil {
		t.Fatal(err)
	}

}

/*
// call read service with token=GITHUB_TOKEN from environment variable
func TestReposSuccess(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	user := "proshik";
	github, err := NewGitHub(token)

	_, err = github.Repos(user)
	if err != nil {
		t.Errorf("Not found not one repository by username=%s", user)
	}
}
 */

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
