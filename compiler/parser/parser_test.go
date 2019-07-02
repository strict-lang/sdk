package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/token"
	"testing"
)

func TestSkipOperator(test *testing.T) {
	entries := map[token.Operator]token.Queue{
		token.AddOperator: {
			token.NewOperatorToken(token.AddOperator, token.Position{Begin: 0, End: 1}, 0),
			token.NewOperatorToken(token.SubOperator, token.Position{Begin: 1, End: 2}, 0),
		},
		token.SubOperator: {
			token.NewOperatorToken(token.SubOperator, token.Position{Begin: 0, End: 1}, 0),
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
	return NewParser("test", tokens, diagnostic.NewRecorder())
}
