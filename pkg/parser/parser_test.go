package parser

import (
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/scope"
	"github.com/BenjaminNitschke/Strict/pkg/token"
	"testing"
)

func TestSkipOperator(test *testing.T) {
	entries := map[token.Operator]token.Queue{
		token.AddOperator: {
			token.NewOperatorToken(token.AddOperator, token.Position{Begin: 0, End: 1}),
			token.NewOperatorToken(token.SubOperator, token.Position{Begin: 1, End: 2}),
		},
		token.SubOperator: {
			token.NewOperatorToken(token.SubOperator, token.Position{Begin: 0, End: 1}),
		},
	}

	for entry, queue := range entries {
		parser := NewTestParser(token.NewQueueReader(queue))
		if err := parser.skipOperator(entry); err != nil {
			test.Errorf("failed expecting operator: %s", err.Error())
		}
	}
}

func NewTestParser(tokens token.Reader) *Parser {
	unit := ast.NewTranslationUnit("test", scope.NewRoot(), []ast.Node{})
	return NewParser(&unit, tokens, diagnostic.NewRecorder())
}
