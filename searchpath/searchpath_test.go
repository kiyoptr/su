package searchpath

import "testing"

func TestLookupInPaths(t *testing.T) {
	found, err := LookupInPaths([]string{
		"~/.?",
		"./?.go",
	}, "searchpath")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(found)
}
