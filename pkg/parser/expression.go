package parser

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

func (parser Parser) ParseUnaryExpression() (ast.Expression, error) {
	switch peek := parser.tokens.Peek(); {
	case peek.IsOperator():
		switch peek.(*token.OperatorToken).Operator {
		}
		break
	}
	return nil, errors.New("Invalid unary expression")
}