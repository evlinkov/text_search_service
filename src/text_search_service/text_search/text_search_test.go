package text_search

import (
	"testing"
)

func TestCorrectnessTextSearchAdd(t *testing.T) {
	textSearch := InitTextSearch(nil)
	response1 := textSearch.AddWord("text1")
	response2 := textSearch.AddWord("text2")
	response3 := textSearch.AddWord("text3")
	response4 := textSearch.AddWord("text1")

	if response1 == response2 || response1 == response3 || response1 != response4 || response2 == response3 {
		t.Fatalf("error while add word")
	}
	t.Logf("success")
}
