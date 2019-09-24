package session

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type AstCache struct {
	files map[string]entry
}

type entry struct {
	unit       *syntaxtree2.TranslationUnit
	lastUpdate int64
}

func (cache *AstCache) find(filePath string) (unit *syntaxtree2.TranslationUnit, found bool) {
	return nil, false
}

func (cache *AstCache) put(filePath string, unit *syntaxtree2.TranslationUnit) {
}

func (cache *AstCache) invalidate() {
}
