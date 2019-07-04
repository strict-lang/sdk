package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

// ParseTypeName is a recursive method that parses type names.
// When calling this method, the types primary name is the value
// of the 'last' token.
// TODO(merlinosayimwen): The scanner currently does not scan
//  bitshift operations, this method however will fail parsing
//  nested generic like `list<list<number>>` because the scanner
//  scans a RightShift operator instead of two GreaterOperators.
func (parser *Parser) ParseTypeName() (ast.TypeName, error) {
	typename := parser.token()
	parser.advance()
	if !token.IsIdentifierToken(typename) {
		return nil, &UnexpectedTokenError{
			Token:    typename,
			Expected: "typename",
		}
	}
	if token.OperatorValue(parser.token()) != token.SmallerOperator {
		return &ast.ConcreteTypeName{
			Name: typename.Value(),
		}, nil
	}
	parser.advance()
	generic, err := parser.ParseTypeName()
	if err != nil {
		return nil, err
	}
	closingOperator := parser.token()
	if token.OperatorValue(closingOperator) != token.GreaterOperator {
		return nil, &UnexpectedTokenError{
			Token:    closingOperator,
			Expected: token.GreaterOperator.String(),
		}
	}
	return &ast.GenericTypeName{
		Name:    typename.Value(),
		Generic: generic,
	}, nil
}
