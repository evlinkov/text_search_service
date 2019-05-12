package text_search

import (
	"fmt"
	"github.com/ljfuyuan/suffixtree"
	"github.com/satori/go.uuid"
	"log"
	"sync"
	"text_search_service/util"
)

const (
	numberOfWorkersForUpdate = 100
)

type TextSearch struct {
	tree  Tree
	words *sync.Map // uuid -> Word
	mutex *sync.RWMutex
	index int

	setOfWords      map[string]uuid.UUID
	mutexSetOfWords *sync.Mutex

	indexToUuid *sync.Map // index -> uuid
	updaters    []*sync.Mutex
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
	textSearch.setOfWords = make(map[string]uuid.UUID)
	textSearch.mutexSetOfWords = &sync.Mutex{}
	textSearch.indexToUuid = &sync.Map{}
	textSearch.updaters = make([]*sync.Mutex, numberOfWorkersForUpdate)
	for i := 0; i < numberOfWorkersForUpdate; i++ {
		textSearch.updaters[i] = &sync.Mutex{}
	}

	for _, word := range words {
		textSearch.words.Store(fmt.Sprintf("%v", word.Uuid), word)
		textSearch.tree.Put(word.Text, word.Index)
		textSearch.indexToUuid.Store(word.Index, word.Uuid)
		if word.Index > textSearch.index {
			textSearch.index = word.Index + 1
		}
		textSearch.setOfWords[word.Text] = word.Uuid
	}
	return textSearch
}

func (textSearch *TextSearch) AddWord(text string) uuid.UUID {
	textSearch.mutexSetOfWords.Lock()
	value, exists := textSearch.setOfWords[text]
	if exists {
		textSearch.mutexSetOfWords.Unlock()
		return value
	}
	word := Word{}
	word.Uuid = util.GenerateUUID()
	word.Text = text
	word.Popularity = 1
	word.Index = textSearch.index
	textSearch.index++
	textSearch.setOfWords[text] = word.Uuid
	textSearch.mutexSetOfWords.Unlock()
	textSearch.addWord(word)
	return word.Uuid
}

func (textSearch *TextSearch) Search(text string) []Word {
	textSearch.mutex.RLock()
	indexes := textSearch.tree.Search(text, -1)
	textSearch.mutex.RUnlock()
	words := make([]Word, 0)
	for _, index := range indexes {
		uuid, ok := textSearch.indexToUuid.Load(index)
		if ok {
			word, ok := textSearch.words.Load(fmt.Sprintf("%v", uuid))
			if ok {
				words = append(words, word.(Word))
				textSearch.updateWord(word.(Word))
			}
		}
	}
	return words
}

func (textSearch *TextSearch) GetWordByUUID(uuid uuid.UUID) (Word, bool) {
	word, ok := textSearch.words.Load(fmt.Sprintf("%v", uuid))
	if ok {
		return word.(Word), ok
	}
	return Word{}, ok
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

func (textSearch *TextSearch) addWord(word Word) {
	textSearch.words.Store(fmt.Sprintf("%v", word.Uuid), word)
	textSearch.indexToUuid.Store(word.Index, word.Uuid)
	defer textSearch.mutex.Unlock()
	textSearch.mutex.Lock()
	textSearch.tree.Put(word.Text, word.Index)
}

func (textSearch *TextSearch) updateWord(word Word) {
	worker := util.GetHash(word.Text) % 100
	defer textSearch.updaters[worker].Unlock()
	textSearch.updaters[worker].Lock()
	object, exists := textSearch.words.Load(fmt.Sprintf("%v", word.Uuid))
	if !exists {
		log.Printf("can not find word in sync map %+v\n", word)
		return
	}
	word.Popularity = object.(Word).Popularity + 1
	textSearch.words.Store(fmt.Sprintf("%v", word.Uuid), word)
}
