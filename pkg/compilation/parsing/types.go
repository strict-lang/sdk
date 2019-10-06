package parsing

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/source"
	"gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/pkg/compilation/token"
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
func (parsing *Parsing) parseTypeName() (syntaxtree.TypeName, error) {
	beginOffset := parsing.offset()
	typeName := parsing.token()
	parsing.advance()
	return parsing.parseTypeNameFromBaseIdentifier(beginOffset, typeName)
}

func (parsing *Parsing) parseTypeNameFromBaseIdentifier(
	beginOffset source.Offset, typeName token.Token) (syntaxtree.TypeName, error) {

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
	concrete := &syntaxtree.ConcreteTypeName{
		Name:         typeName.Value(),
		NodePosition: parsing.createPosition(beginOffset),
	}
	if operator == token.LeftBracketOperator {
		return parsing.parseListTypeName(beginOffset, concrete)
	}
	return concrete, nil
}

func (parsing *Parsing) parseGenericTypeName(
	beginOffset source.Offset, base string) (syntaxtree.TypeName, error) {

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
	return &syntaxtree.GenericTypeName{
		Name:         base,
		Generic:      generic,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseListTypeName(
	beginOffset source.Offset, base syntaxtree.TypeName) (syntaxtree.TypeName, error) {

	if err := parsing.skipOperator(token.LeftBracketOperator); err != nil {
		return nil, err
	}
	if err := parsing.skipOperator(token.RightBracketOperator); err != nil {
		return nil, &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: "], end of list name",
		}
	}
	typeName := &syntaxtree.ListTypeName{
		ElementTypeName: base,
		NodePosition:    parsing.createPosition(beginOffset),
	}
	if token.HasOperatorValue(parsing.token(), token.LeftBracketOperator) {
		return parsing.parseListTypeName(beginOffset, typeName)
	}
	return typeName, nil
}
