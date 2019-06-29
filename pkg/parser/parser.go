package parser

import (
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

type Parser struct {
	tokens   token.Reader
	recorder diagnostics.Recorder
}

func NewParser(tokens token.Reader) {

}
