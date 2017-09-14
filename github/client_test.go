package github

import (
	"net/http/httptest"
	"net/http"
	"net/url"
	"github.com/google/go-github/github"
	"testing"
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


func TestNewClient(t *testing.T) {
	token := "token"

	_, err := NewClient(token)
	if err != nil {
		t.Fatal(err)
	}

}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
