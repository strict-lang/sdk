package session

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

type AstCache struct {
	files map[string]entry
}

type entry struct {
	unit       *ast.TranslationUnit
	lastUpdate int64
}

func (cache *AstCache) find(filePath string) (unit *ast.TranslationUnit, found bool) {
	return nil, false
}

func (cache *AstCache) put(filePath string, unit *ast.TranslationUnit) {
}

func (cache *AstCache) invalidate() {
}
