package parsing

import (
	diagnostic2 "gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

// skipOperator skips the next keyword if it the passed operator, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipOperator(operator token2.Operator) error {
	if err := parsing.expectOperator(operator); err != nil {
		return err
	}
	parsing.advance()
	return nil
}

// skipKeyword skips the next keyword if it the passed keyword, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipKeyword(keyword token2.Keyword) error {
	if err := parsing.expectKeyword(keyword); err != nil {
		return err
	}
	parsing.advance()
	return nil
}

// expectOperator peeks the next token and expects it to be the passed operator,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectOperator(expected token2.Operator) error {
	if token2.OperatorValue(parsing.token()) != expected {
		return &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: expected.String(),
		}
	}
	return nil
}

// expectKeyword peeks the next token and expects it to be the passed keyword,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectKeyword(expected token2.Keyword) error {
	if token2.KeywordValue(parsing.token()) != expected {
		return &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: expected.String(),
		}
	}
	return nil
}

// expectAnyIdentifier expects the next token to be an identifier,
// without regards to its value and returns an error if it fails.
func (parsing *Parsing) expectAnyIdentifier() (*syntaxtree2.Identifier, error) {
	current := parsing.token()
	if !token2.IsIdentifierToken(current) {
		return nil, &UnexpectedTokenError{
			Token:    current,
			Expected: "any identifier",
		}
	}
	return &syntaxtree2.Identifier{
		Value:        current.Value(),
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) isLookingAtKeyword(keyword token2.Keyword) bool {
	return token2.HasKeywordValue(parsing.peek(), keyword)
}

func (parsing *Parsing) isLookingAtOperator(operator token2.Operator) bool {
	return token2.HasOperatorValue(parsing.peek(), operator)
}

func (parsing *Parsing) createInvalidStatement(beginOffset source2.Offset, err error) syntaxtree2.Node {
	parsing.reportError(err, parsing.createPosition(beginOffset))
	return &syntaxtree2.InvalidStatement{
		NodePosition: parsing.createPosition(beginOffset),
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
func (parsing *Parsing) reportError(err error, position syntaxtree2.Position) {
	parsing.recorder.Record(diagnostic2.RecordedEntry{
		Kind:     &diagnostic2.Error,
		Stage:    &diagnostic2.SyntacticalAnalysis,
		Message:  err.Error(),
		Position: position,
	})
}

func (parsing *Parsing) createTokenPosition() syntaxtree2.Position {
	return parsing.token().Position()
}

func (parsing *Parsing) createPosition(beginOffset source2.Offset) syntaxtree2.Position {
	return &offsetPosition{begin: beginOffset, end: parsing.offset()}
}

type offsetPosition struct {
	begin source2.Offset
	end   source2.Offset
}

// Begin returns the offset to the position begin.
func (position offsetPosition) Begin() source2.Offset {
	return position.begin
}

// End returns the offset to the positions end.
func (position offsetPosition) End() source2.Offset {
	return position.end
}
