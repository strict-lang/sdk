package parsing

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

func (parsing *Parsing) couldBeLookingAtTypeName() bool {
	if !token2.IsIdentifierToken(parsing.token()) {
		return false
	}
	if token2.IsIdentifierToken(parsing.peek()) {
		return true
	}
	peek := parsing.peek()
	return token2.HasOperatorValue(peek, token2.SmallerOperator) ||
		token2.HasOperatorValue(peek, token2.LeftBracketOperator)
}

// parseTypeName is a recursive method that parses type names. When calling
// this method, the types primary name is the value of the 'last' token.
func (parsing *Parsing) parseTypeName() (syntaxtree2.TypeName, error) {
	beginOffset := parsing.offset()
	typeName := parsing.token()
	parsing.advance()
	if !token2.IsIdentifierToken(typeName) {
		return nil, &UnexpectedTokenError{
			Token:    typeName,
			Expected: "TypeName",
		}
	}
	operator := token2.OperatorValue(parsing.token())
	if operator == token2.SmallerOperator {
		return parsing.parseGenericTypeName(beginOffset, typeName.Value())
	}
	concrete := &syntaxtree2.ConcreteTypeName{
		Name:         typeName.Value(),
		NodePosition: parsing.createPosition(beginOffset),
	}
	if operator == token2.LeftBracketOperator {
		return parsing.parseListTypeName(beginOffset, concrete)
	}
	return concrete, nil
}

func (parsing *Parsing) parseGenericTypeName(
	beginOffset source2.Offset, base string) (syntaxtree2.TypeName, error) {

	if err := parsing.skipOperator(token2.SmallerOperator); err != nil {
		return nil, err
	}
	generic, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	closingOperator := parsing.token()
	if token2.OperatorValue(closingOperator) != token2.GreaterOperator {
		return nil, &UnexpectedTokenError{
			Token:    closingOperator,
			Expected: token2.GreaterOperator.String(),
		}
	}
	parsing.advance()
	return &syntaxtree2.GenericTypeName{
		Name:         base,
		Generic:      generic,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseListTypeName(
	beginOffset source2.Offset, base syntaxtree2.TypeName) (syntaxtree2.TypeName, error) {

	if err := parsing.skipOperator(token2.LeftBracketOperator); err != nil {
		return nil, err
	}
	if err := parsing.skipOperator(token2.RightBracketOperator); err != nil {
		return nil, err
	}
	typeName := &syntaxtree2.ListTypeName{
		ElementTypeName: base,
		NodePosition:    parsing.createPosition(beginOffset),
	}
	if token2.HasOperatorValue(parsing.token(), token2.LeftBracketOperator) {
		return parsing.parseListTypeName(beginOffset, typeName)
	}
	return typeName, nil
}
