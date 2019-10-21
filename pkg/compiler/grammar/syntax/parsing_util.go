package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
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

// expectAnyIdentifier expects the next token to be an identifier,
// without regards to its value and returns an error if it fails.
func (parsing *Parsing) expectAnyIdentifier() *tree.Identifier {
	current := parsing.token()
	if !token.IsIdentifierToken(current) {
		parsing.throwError(&UnexpectedTokenError{
			Token:    current,
			Expected: "any identifier",
		})
	}
	return &tree.Identifier{
		Value:  current.Value(),
		Region: parsing.createRegionFromCurrentToken(),
	}
}

func (parsing *Parsing) isLookingAtKeyword(keyword token.Keyword) bool {
	return token.HasKeywordValue(parsing.peek(), keyword)
}

func (parsing *Parsing) isLookingAtOperator(operator token.Operator) bool {
	return token.HasOperatorValue(parsing.peek(), operator)
}

func (parsing *Parsing) createInvalidStatement(beginOffset input.Offset, err error) tree.Statement {
	parsing.reportError(err, parsing.createRegion(beginOffset))
	return &tree.InvalidStatement{
		Region: parsing.createRegion(beginOffset),
	}
}

// skipEndOfStatement skips the next token if it is an EndOfStatement token.
func (parsing *Parsing) skipEndOfStatement() {
	// Do not report the missing end of statement.
	parsing.advance()
	parsing.isAtBeginOfStatement = true
}

// reportError reports an error to the diagnostics bag, starting at the
// passed position and ending at the parsers current position.
func (parsing *Parsing) reportError(err error, position input.Region) {
	parsing.recorder.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SyntacticalAnalysis,
		Message:  err.Error(),
		Position: position,
	})
}

func (parsing *Parsing) createRegionFromCurrentToken() input.Region {
	return parsing.createRegionFromToken(parsing.token())
}

func (parsing *Parsing) createRegionFromToken(token token.Token) input.Region {
	tokenPosition := token.Position()
	return input.CreateRegion(tokenPosition.Begin(), tokenPosition.End())
}

func (parsing *Parsing) createRegion(beginOffset input.Offset) input.Region {
	return input.CreateRegion(beginOffset, parsing.offset())
}

func (parsing *Parsing) throwError(err error) {}
