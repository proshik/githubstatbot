package github

import (
	"net/http"
	"encoding/json"
	_"fmt"
	"fmt"
)

type Repos struct {
	Repo []*Repo
}

type Repo struct {
	Id   int `json:"id"`
	Name string `json:"name"`
}

type Language struct {
}

func allRepos(profile string) (map[string]int, error) {
	resp, err := http.Get("https://api.github.com/users/" + profile + "/repos")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data []*Repo
	json.NewDecoder(resp.Body).Decode(&data)

	result := make(map[string]int, 0)

	for _, d := range data {
		resp1, err := http.Get("https://api.github.com/repos/" + profile + "/" + d.Name + " /languages")
		if err != nil {
			return nil, err
		}
		defer resp1.Body.Close()

		var langs map[string]int
		json.NewDecoder(resp1.Body).Decode(&langs)



		result = append(result, langs[])
	}

	return make(map[string]int, 0), nil

}
