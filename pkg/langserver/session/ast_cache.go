package session

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type AstCache struct {
	files map[string]entry
}

type entry struct {
	unit       *syntaxtree.TranslationUnit
	lastUpdate int64
}

func (cache *AstCache) find(filePath string) (unit *syntaxtree.TranslationUnit, found bool) {
	return nil, false
}

func (cache *AstCache) put(filePath string, unit *syntaxtree.TranslationUnit) {
}

func (cache *AstCache) invalidate() {
}
