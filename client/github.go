package client

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"context"
	_"log"
	_"sort"
)

type GitHub struct {
	Client *github.Client
}

func NewGithub(token string) (*GitHub, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)

	return &GitHub{c}, nil
}

func (github *GitHub) Repos(user string) ([]*github.Repository, error) {

	ctx := context.Background()

	repos, _, err := github.Client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (github *GitHub) Language(user string, repoName string) (map[string]int, error) {
	ctx := context.Background()

	lang, _, err := github.Client.Repositories.ListLanguages(ctx, user, repoName)
	if err != nil {
		return nil, err
	}

	return lang, nil
}

//func (github *GitHub) Languages(username string) ([]*Languages, error) {
//
//	ctx := context.Background()
//
//	repos, _, err := github.Client.Repositories.List(ctx, username, nil)
//	if err != nil {
//		return nil, err
//	}
//	//calculate total bytes by language
//	language := make(map[string]int)
//	for _, repo := range repos {
//		lang, _, err := github.Client.Repositories.ListLanguages(ctx, username, *repo.Name)
//		if err != nil {
//			log.Printf("Error on request language statistic for repository=%s with err=%v", repo.Name, err)
//		}
//
//		for k, v := range lang {
//			language[k] = language[k] + v
//		}
//	}
//
//	return calcPercentages(language), nil
//}
//
//func calcPercentages(languages map[string]int) []*Languages {
//
//	result := make([]*Languages, 0)
//
//	var totalSum float32
//	for _, v := range languages {
//		totalSum += float32(v)
//	}
//
//	for key, value := range languages {
//		percent := float32(value) * (float32(100) / totalSum)
//		result = append(result, &Languages{key, round(percent, 0.1)})
//	}
//
//	sort.Slice(result, func(i, j int) bool { return result[i].Percentage > result[j].Percentage })
//
//	return result
//}
//
//func round(x, unit float32) float32 {
//	if x > 0 {
//		return float32(int32(x/unit+0.5)) * unit
//	}
//	return float32(int32(x/unit-0.5)) * unit
//}
