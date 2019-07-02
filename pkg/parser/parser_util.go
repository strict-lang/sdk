package parser

import "github.com/BenjaminNitschke/Strict/pkg/token"

// skipOperator skips the next keyword if it the passed operator, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) skipOperator(operator token.Operator) error {
	if err := parser.expectOperator(operator); err != nil {
		return err
	}
	parser.tokens.Pull()
	return nil
}

// skipKeyword skips the next keyword if it the passed keyword, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) skipKeyword(keyword token.Keyword) (bool, error) {
	if err := parser.expectKeyword(keyword); err != nil {
		return false, err
	}
	parser.tokens.Pull()
	return true, nil
}

// expectOperator peeks the next token and expects it to be the passed operator,
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) expectOperator(expected token.Operator) error {
	peek := parser.tokens.Peek()
	if !token.IsOperatorToken(peek) || peek.(*token.OperatorToken).Operator != expected {
		return &UnexpectedTokenError{
			Token:    peek,
			Expected: expected.String(),
		}
	}
	return nil
}

// expectKeyword peeks the next token and expects it to be the passed keyword,
// otherwise an UnexpectedTokenError is returned.
func (parser *Parser) expectKeyword(expected token.Keyword) error {
	peek := parser.tokens.Peek()
	if !token.IsKeywordToken(peek) || peek.(*token.KeywordToken).Keyword != expected {
		return &UnexpectedTokenError{
			Token:    peek,
			Expected: expected.String(),
		}
	}
	return nil
}

func (parser *Parser) expectAnyIdentifier() error {
	peek := parser.tokens.Peek()
	if peek.Name() != token.IdentifierTokenName {
		return &UnexpectedTokenError{
			Token: peek,
			Expected: "any identifier",
		}
	}
	return nil
}

func (parser *Parser) isLookingAtKeyword(keyword token.Keyword) bool {
	peek := parser.tokens.Peek()
	if !token.IsKeywordToken(peek) {
		return false
	}
	return peek.(*token.KeywordToken).Keyword == keyword
}

func (parser *Parser) isLookingAtOperator(operator token.Operator) bool {
	peek := parser.tokens.Peek()
	if !token.IsOperatorToken(peek) {
		return parser.isLookingAtOperatorKeyword(operator)
	}
	return peek.(*token.OperatorToken).Operator == operator
}

func (parser *Parser) isLookingAtOperatorKeyword(operator token.Operator) bool {
	peek := parser.tokens.Peek()
	if !token.IsKeywordToken(peek) {
		return false
	}
	keyword := peek.(*token.KeywordToken)
	return keyword.AsOperator() == operator
}