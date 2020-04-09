package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

// skipOperator skips the next keyword if it the passed operator, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipOperator(operator token.Operator) {
	if err := parsing.expectOperator(operator); err != nil {
		parsing.throwError(err)
	}
	parsing.advance()
}

// skipKeyword skips the next keyword if it the passed keyword, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipKeyword(keyword token.Keyword) {
	if err := parsing.expectKeyword(keyword); err != nil {
		parsing.throwError(err)
	}
	parsing.advance()
}

// expectOperator peeks the next token and expects it to be the passed operator,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectOperator(expected token.Operator) *diagnostic.RichError {
	if token.OperatorValue(parsing.token()) != expected {
		return newInvalidOperatorError(parsing.token(), expected)
	}
	return nil
}

// expectKeyword peeks the next token and expects it to be the passed keyword,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectKeyword(expected token.Keyword) *diagnostic.RichError {
	if token.KeywordValue(parsing.token()) != expected {
		return newInvalidKeywordError(parsing.token(), expected)
	}
	return nil
}

// expectAnyIdentifier expects the next token to be an identifier,
// without regards to its value and returns an error if it fails.
func (parsing *Parsing) expectAnyIdentifier() *tree.Identifier {
	parsing.beginStructure(tree.IdentifierNodeKind)
	current := parsing.token()
	if !token.IsIdentifierToken(current) {
		parsing.throwError(newNoIdentifierError(current))
	}
	return &tree.Identifier{
		Value:  current.Value(),
		Region: parsing.completeStructure(tree.IdentifierNodeKind),
	}
}

func (parsing *Parsing) isLookingAtKeyword(keyword token.Keyword) bool {
	return token.HasKeywordValue(parsing.peek(), keyword)
}

func (parsing *Parsing) isLookingAtOperator(operator token.Operator) bool {
	return token.HasOperatorValue(parsing.peek(), operator)
}

func (parsing *Parsing) completeInvalidStructure(err error) tree.Statement {
	region := parsing.completeStructure(tree.WildcardNodeKind)
	parsing.reportError(newInvalidStructureError(), region)
	return &tree.InvalidStatement{
		Region: region,
	}
}
