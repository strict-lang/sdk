package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

func (parsing *Parsing) couldBeLookingAtTypeName() bool {
	if !token.IsIdentifierToken(parsing.token()) {
		return false
	}
	if token.IsIdentifierToken(parsing.peek()) {
		return true
	}
	return token.HasOperatorValue(parsing.peek(), token.SmallerOperator)
}

// parseTypeName is a recursive method that parses type names. When calling
// this method, the types primary name is the value of the 'last' token.
func (parsing *Parsing) parseTypeName() (ast.TypeName, error) {
	beginOffset := parsing.offset()
	typeName := parsing.token()
	parsing.advance()
	if !token.IsIdentifierToken(typeName) {
		return nil, &UnexpectedTokenError{
			Token:    typeName,
			Expected: "TypeName",
		}
	}
	if token.OperatorValue(parsing.token()) != token.SmallerOperator {
		return &ast.ConcreteTypeName{
			Name: typeName.Value(),
			NodePosition: parsing.createPosition(beginOffset)	,
		}, nil
	}
	parsing.advance()
	return parsing.parseGenericTypeName(beginOffset, typeName.Value())
}

func (parsing *Parsing) parseGenericTypeName(
	beginOffset source.Offset, base string) (ast.TypeName, error) {

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
		Name:    base,
		Generic: generic,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}