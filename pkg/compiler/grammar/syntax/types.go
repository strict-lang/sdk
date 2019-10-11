package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

func (parsing *Parsing) couldBeLookingAtTypeName() bool {
	if !token.IsIdentifierToken(parsing.token()) {
		return false
	}
	peek := parsing.peek()
	if token.IsIdentifierToken(peek) {
		return true
	}
	return token.HasOperatorValue(peek, token.SmallerOperator) ||
		token.HasOperatorValue(peek, token.LeftBracketOperator)
}

// parseTypeName is a recursive method that parses type names. When calling
// this method, the types primary name is the value of the 'last' token.
func (parsing *Parsing) parseTypeName() (tree.TypeName, error) {
	beginOffset := parsing.offset()
	typeName := parsing.token()
	parsing.advance()
	return parsing.parseTypeNameFromBaseIdentifier(beginOffset, typeName)
}

func (parsing *Parsing) parseTypeNameFromBaseIdentifier(
	beginOffset input.Offset, typeName token.Token) (tree.TypeName, error) {

	if !token.IsIdentifierToken(typeName) {
		return nil, &UnexpectedTokenError{
			Token:    typeName,
			Expected: "TypeName",
		}
	}
	operator := token.OperatorValue(parsing.token())
	if operator == token.SmallerOperator {
		return parsing.parseGenericTypeName(beginOffset, typeName.Value())
	}
	concrete := &tree.ConcreteTypeName{
		Name:         typeName.Value(),
		Region: parsing.createRegion(beginOffset),
	}
	if operator == token.LeftBracketOperator {
		return parsing.parseListTypeName(beginOffset, concrete)
	}
	return concrete, nil
}

func (parsing *Parsing) parseGenericTypeName(
	beginOffset input.Offset, base string) (tree.TypeName, error) {

	if err := parsing.skipOperator(token.SmallerOperator); err != nil {
		return nil, err
	}
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
	return &tree.GenericTypeName{
		Name:         base,
		Generic:      generic,
		Region: parsing.createRegion(beginOffset),
	}, nil
}

func (parsing *Parsing) parseListTypeName(
	beginOffset input.Offset, base tree.TypeName) (tree.TypeName, error) {

	if err := parsing.skipOperator(token.LeftBracketOperator); err != nil {
		return nil, err
	}
	if err := parsing.skipOperator(token.RightBracketOperator); err != nil {
		return nil, &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: "], end of list name",
		}
	}
	typeName := &tree.ListTypeName{
		ElementTypeName: base,
		Region:    parsing.createRegion(beginOffset),
	}
	if token.HasOperatorValue(parsing.token(), token.LeftBracketOperator) {
		return parsing.parseListTypeName(beginOffset, typeName)
	}
	return typeName, nil
}
