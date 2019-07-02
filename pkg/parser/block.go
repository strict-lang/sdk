package parser

import (
	"github.com/BenjaminNitschke/Strict/pkg/scope"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

type Block struct {
	Indent token.Indent
	Scope  *scope.Scope
	Parent *Block
}

