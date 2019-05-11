package text_search

import (
	"github.com/ljfuyuan/suffixtree"
	"github.com/satori/go.uuid"
	"sync"
)

type TextSearch struct {
	tree  Tree
	words *sync.Map
}

type Word struct {
	text       string    `json:"text"`
	uuid       uuid.UUID `json:"uuid"`
	popularity int64     `json:"popularity"`
	index      int       `json:"index"`
}

func InitTextSearch(words []Word) *TextSearch {
	textSearch := &TextSearch{}
	textSearch.tree = suffixtree.NewGeneralizedSuffixTree()
	textSearch.words = &sync.Map{}
	for _, word := range words {
		textSearch.words.Store(word.uuid, word)
		textSearch.tree.Put(word.text, word.index)
	}
	return textSearch
}

func (textSearch *TextSearch) GetAllWords() []Word {
	words := make([]Word, 0)
	textSearch.words.Range(func(key interface{}, value interface{}) bool {
		words = append(words, value.(Word))
		return true
	})
	return words
}
