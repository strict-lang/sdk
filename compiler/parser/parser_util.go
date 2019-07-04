package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/token"
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

func (parser *Parser) expectAnyIdentifier() (ast.Identifier, error) {
	current := parser.token()
	if !token.IsIdentifierToken(current) {
		return ast.Identifier{}, &UnexpectedTokenError{
			Token:    parser.token(),
			Expected: "any identifier",
		}
	}
	return ast.Identifier{
		Value: current.Value(),
	}, nil
}

func (parser *Parser) isLookingAtKeyword(keyword token.Keyword) bool {
	return token.HasKeywordValue(parser.peek(), keyword)
}

func (parser *Parser) isLookingAtOperator(operator token.Operator) bool {
	return token.HasOperatorValue(parser.peek(), operator)
}

func (parser *Parser) createInvalidStatement(err error) ast.Node {
	parser.reportError(err)
	return &ast.InvalidStatement{}
}

// skipEndOfStatement skips the next token if it is an EndOfStatement token.
func (parser *Parser) skipEndOfStatement() {
	// Do not report the missing end of statement.
	parser.advance()
}
