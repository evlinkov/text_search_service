package text_search

type Tree interface {
	Search(word string, numElements int) []int
	Put(key string, index int)
}
