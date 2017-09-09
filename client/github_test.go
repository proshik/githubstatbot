package client

import (
	"testing"
	"os"
	_"fmt"
	_"sort"
)

func TestNewGithubClient(t *testing.T){

	token := os.Getenv("GITHUB_TOKEN")

	_, err := NewGithub(token)
	if err != nil {
		t.Fatal(err)
	}

}

//func TestAllRepos(t *testing.T) {
//	token := "token"
//
//	auth := &GHAuth{token}
//
//	result, err := auth.("proshik")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(result) <= 0 {
//		t.Fatal("Lenght of map is empty")
//	}
//
//	fmt.Printf("Result: %s\n", result)
//}
//
//func TestCalcPercentages(t *testing.T) {
//	languages := map[string]int{
//		"Java": 199992,
//		"Go":   16172,
//		"HTML": 13579}
//
//	expected := []*Languages{
//		{"Java", 87.1},
//		{"Go", 7},
//		{"HTML", 5.9},
//	}
//	sort.Slice(expected, func(i, j int) bool {
//		return expected[i].Percentage > expected[j].Percentage
//	})
//
//	result := calcPercentages(languages)
//
//	for i := 0; i < len(result); i++ {
//		if *result[i] != *expected[i] {
//			t.Errorf("CalcPercentages = expected=%v, actual=%v", expected[i], result[i])
//		}
//	}
//	//fmt.Sprintf("%.1f", result[0])
//}
