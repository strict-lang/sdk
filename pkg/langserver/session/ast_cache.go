package session

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/grammar/syntax/tree"
)

type AstCache struct {
	files map[string]entry
}

type entry struct {
	unit       *tree.TranslationUnit
	lastUpdate int64
}

func (cache *AstCache) find(filePath string) (unit *tree.TranslationUnit, found bool) {
	return nil, false
}

func (cache *AstCache) put(filePath string, unit *tree.TranslationUnit) {
}

func (cache *AstCache) invalidate() {
}
