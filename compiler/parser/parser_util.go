package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

// skipOperator skips the next keyword if it the passed operator, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) skipOperator(operator token.Operator) error {
	if err := parser.expectOperator(operator); err != nil {
		return err
	}
	parser.advance()
	return nil
}

// skipKeyword skips the next keyword if it the passed keyword, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) skipKeyword(keyword token.Keyword) error {
	if err := parser.expectKeyword(keyword); err != nil {
		return err
	}
	parser.advance()
	return nil
}

// expectOperator peeks the next token and expects it to be the passed operator,
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) expectOperator(expected token.Operator) error {
	if token.OperatorValue(parser.token()) != expected {
		return &UnexpectedTokenError{
			Token:    parser.token(),
			Expected: expected.String(),
		}
	}
	return nil
}

// expectKeyword peeks the next token and expects it to be the passed keyword,
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) expectKeyword(expected token.Keyword) error {
	if token.KeywordValue(parser.token()) != expected {
		return &UnexpectedTokenError{
			Token:    parser.token(),
			Expected: expected.String(),
		}
	}
	return nil
}

func (parser *Parser) expectAnyIdentifier() (*ast.Identifier, error) {
	current := parser.token()
	if !token.IsIdentifierToken(current) {
		return nil, &UnexpectedTokenError{
			Token:    current,
			Expected: "any identifier",
		}
	}
	return &ast.Identifier{
		Value:        current.Value(),
		NodePosition: parser.createTokenPosition(),
	}, nil
}

func (parser *Parser) isLookingAtKeyword(keyword token.Keyword) bool {
	return token.HasKeywordValue(parser.peek(), keyword)
}

func (parser *Parser) isLookingAtOperator(operator token.Operator) bool {
	return token.HasOperatorValue(parser.peek(), operator)
}

func (parser *Parser) createInvalidStatement(beginOffset source.Offset, err error) ast.Node {
	parser.reportError(err)
	return &ast.InvalidStatement{
		NodePosition: parser.createPosition(beginOffset),
	}
}

// skipEndOfStatement skips the next token if it is an EndOfStatement token.
func (parser *Parser) skipEndOfStatement() {
	// Do not report the missing end of statement.
	parser.advance()
}
