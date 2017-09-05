package github

import (
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type GHAuth struct {
	Token string
}

type LangStatistic struct {
	Language   string
	Percentage float32
}

type Repo struct {
	Id   int `json:"id"`
	Name string `json:"name"`
}

var client = http.Client{}

func (auth *GHAuth) LangStatistic(profile string) ([]*LangStatistic, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", profile)

	resp, err := get(auth, url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var repos []*Repo
	json.NewDecoder(resp.Body).Decode(&repos)

	result := make(map[string]int)
	for _, repo := range repos {
		lang, err := getLanguages(profile, repo.Name, auth)
		if err != nil {
			log.Printf("Error on request language statistic for repository=%s with err=%v", repo.Name, err)
		}

		for k, v := range lang {
			result[k] = result[k] + v
		}
	}

	return calcPercentages(result), nil
}

func getLanguages(profile string, repoName string, auth *GHAuth) (map[string]int, error) {

	url := fmt.Sprintf("http://api.github.com/repos/%s/%s/languages", profile, repoName)

	resp, err := get(auth, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var langMap map[string]int
	json.NewDecoder(resp.Body).Decode(&langMap)

	return langMap, nil
}

func get(auth *GHAuth, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+auth.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func calcPercentages(languages map[string]int) []*LangStatistic {

	result := make([]*LangStatistic, 0)

	var totalSum float32

	for _, v := range languages {
		totalSum += float32(v)
	}

	for key, value := range languages {
		percent := float32(value) * (float32(100) / totalSum)
		result = append(result, &LangStatistic{key, round(percent, 0.1)})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Percentage > result[j].Percentage })

	return result
}

func round(x, unit float32) float32 {
	if x > 0 {
		return float32(int32(x/unit+0.5)) * unit
	}
	return float32(int32(x/unit-0.5)) * unit
}
