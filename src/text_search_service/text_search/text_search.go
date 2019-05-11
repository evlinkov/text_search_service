package text_search

import "github.com/ljfuyuan/suffixtree"

const ()

type TextSearch struct {
	tree Tree
}

func InitTextSearch() *TextSearch {
	textSearch := &TextSearch{}
	textSearch.tree = suffixtree.NewGeneralizedSuffixTree()
	return textSearch
}
