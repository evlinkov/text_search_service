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

func TestCorrectnessTextSearchGet(t *testing.T) {
	textSearch := InitTextSearch(nil)
	response := textSearch.AddWord("text1")

	word, ok := textSearch.GetWordByUUID(response)
	if !ok || word.Text != "text1" {
		t.Fatalf("error get word")
	}
	t.Logf("success")
}

func TestCorrectnessTextSearchMethodSearch(t *testing.T) {
	textSearch := InitTextSearch(nil)
	textSearch.AddWord("text1")
	textSearch.AddWord("text2")

	words := textSearch.Search("ext1")
	if len(words) != 1 || words[0].Text != "text1" || words[0].Popularity != 1 {
		t.Fatalf("error search words")
	}
	words = textSearch.Search("ec")
	if len(words) != 0 {
		t.Fatalf("error search words")
	}
	words = textSearch.Search("text")
	if len(words) != 2 || words[0].Text != "text1" || words[0].Popularity != 2 || words[1].Text != "text2" || words[1].Popularity != 1 {
		t.Fatalf("error search words")
	}

	t.Logf("success")
}
