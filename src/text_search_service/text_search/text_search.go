package text_search

import (
	"fmt"
	"github.com/ljfuyuan/suffixtree"
	"github.com/satori/go.uuid"
	"sync"
)

type TextSearch struct {
	tree  Tree
	words *sync.Map
	mutex *sync.RWMutex
}

type Word struct {
	Text       string    `json:"text"`
	Uuid       uuid.UUID `json:"uuid"`
	Popularity int64     `json:"popularity"`
	Index      int       `json:"index"`
}

func InitTextSearch(words []Word) *TextSearch {
	textSearch := &TextSearch{}
	textSearch.tree = suffixtree.NewGeneralizedSuffixTree()
	textSearch.words = &sync.Map{}
	textSearch.mutex = &sync.RWMutex{}
	for _, word := range words {
		textSearch.words.Store(fmt.Sprintf("%v", word.Uuid), word)
		textSearch.tree.Put(word.Text, word.Index)
	}
	return textSearch
}

func (textSearch *TextSearch) GetAllWords() []Word {
	defer textSearch.mutex.RUnlock()
	textSearch.mutex.RLock()
	words := make([]Word, 0)
	textSearch.words.Range(func(key interface{}, value interface{}) bool {
		words = append(words, value.(Word))
		return true
	})
	return words
}
