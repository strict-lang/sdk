package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

// ParseTypeName is a recursive method that parses type names.
// When calling this method, the types primary name is the value
// of the 'last' token.
// TODO(merlinosayimwen): The scanning currently does not scanning
//  bitshift operations, this method however will fail parsing
//  nested generic like `list<list<number>>` because the scanning
//  scans a RightShift operator instead of two GreaterOperators.
func (parsing *Parsing) parseTypeName() (ast.TypeName, error) {
	typename := parsing.token()
	parsing.advance()
	if !token.IsIdentifierToken(typename) {
		return nil, &UnexpectedTokenError{
			Token:    typename,
			Expected: "typename",
		}
	}
	if token.OperatorValue(parsing.token()) != token.SmallerOperator {
		return &ast.ConcreteTypeName{
			Name: typename.Value(),
		}, nil
	}
	parsing.advance()
	generic, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	closingOperator := parsing.token()
	if token.OperatorValue(closingOperator) != token.GreaterOperator {
		return nil, &UnexpectedTokenError{
			Token:    closingOperator,
			Expected: token.GreaterOperator.String(),
		}
	}
	parsing.advance()
	return &ast.GenericTypeName{
		Name:    typename.Value(),
		Generic: generic,
	}, nil
}
