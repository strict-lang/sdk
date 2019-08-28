package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

// skipOperator skips the next keyword if it the passed operator, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipOperator(operator token.Operator) error {
	if err := parsing.expectOperator(operator); err != nil {
		return err
	}
	parsing.advance()
	return nil
}

// skipKeyword skips the next keyword if it the passed keyword, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipKeyword(keyword token.Keyword) error {
	if err := parsing.expectKeyword(keyword); err != nil {
		return err
	}
	parsing.advance()
	return nil
}

// expectOperator peeks the next token and expects it to be the passed operator,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectOperator(expected token.Operator) error {
	if token.OperatorValue(parsing.token()) != expected {
		return &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: expected.String(),
		}
	}
	return nil
}

// expectKeyword peeks the next token and expects it to be the passed keyword,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectKeyword(expected token.Keyword) error {
	if token.KeywordValue(parsing.token()) != expected {
		return &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: expected.String(),
		}
	}
	return nil
}

func (parsing *Parsing) expectAnyIdentifier() (*ast.Identifier, error) {
	current := parsing.token()
	if !token.IsIdentifierToken(current) {
		return nil, &UnexpectedTokenError{
			Token:    current,
			Expected: "any identifier",
		}
	}
	return &ast.Identifier{
		Value:        current.Value(),
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) isLookingAtKeyword(keyword token.Keyword) bool {
	return token.HasKeywordValue(parsing.peek(), keyword)
}

func (parsing *Parsing) isLookingAtOperator(operator token.Operator) bool {
	return token.HasOperatorValue(parsing.peek(), operator)
}

func (parsing *Parsing) createInvalidStatement(beginOffset source.Offset, err error) ast.Node {
	parsing.reportError(err, parsing.createPosition(beginOffset))
	return &ast.InvalidStatement{
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// skipEndOfStatement skips the next token if it is an EndOfStatement token.
func (parsing *Parsing) skipEndOfStatement() {
	// Do not report the missing end of statement.
	parsing.advance()
}
