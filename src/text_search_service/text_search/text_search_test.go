package text_search

import (
	"fmt"
	"github.com/satori/go.uuid"
	"sync"
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

func TestCorrectnessParallelSaveAndGet(t *testing.T) {
	type TestCase struct {
		id   int
		uuid uuid.UUID
	}
	textSearch := InitTextSearch(nil)
	testCases := make([]TestCase, 10)
	wg := &sync.WaitGroup{}
	for i := 0; i < len(testCases); i++ {
		testCases[i].id = i
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			testCases[i].uuid = textSearch.AddWord(fmt.Sprintf("text_%d", i))
		}(i)
	}
	wg.Wait()
	for _, testCase := range testCases {
		wg.Add(1)
		go func(testCase TestCase) {
			defer wg.Done()
			word, ok := textSearch.GetWordByUUID(testCase.uuid)
			if !ok {
				t.Fatalf("error get words")
			}
			if word.Uuid != testCase.uuid || word.Text != fmt.Sprintf("text_%d", testCase.id) {
				t.Fatalf("error get words")
			}
		}(testCase)
	}
	wg.Wait()
	words := textSearch.Search("text")
	if len(words) != 10 {
		t.Fatalf("error search words")
	}

	t.Logf("success")
}

func TestCorrectnessParallelSaveAndSearch(t *testing.T) {
	textSearch := InitTextSearch(nil)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			textSearch.AddWord(fmt.Sprintf("text_%d", i))
		}(i)
	}
	wg.Wait()
	for i := 0; i < 99; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			textSearch.Search("text")
		} ()
	}
	wg.Wait()
	words := textSearch.Search("text")
	if len(words) != 10 {
		t.Fatalf("error search words")
	}
	for _, word := range words {
		if word.Popularity != 100 {
			t.Fatalf("error score popularity")
		}
	}
	t.Logf("success")
}