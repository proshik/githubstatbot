package github

import "testing"

func TestAllRepos(t *testing.T){

	result, err := allRepos("proshik")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) <= 0 {
		t.Fatal("Lenght of map is empty")
	}
}