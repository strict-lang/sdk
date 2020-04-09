package syntax

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"log"
)

func newInvalidStructureError() *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.InvalidStatementError{Kind: tree.UnknownNodeKind},
	}
}

// skipEndOfStatement skips the next token if it is an EndOfStatement token.
func (parsing *Parsing) skipEndOfStatement() {
	// Do not report the missing end of statement.
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
	}
	parsing.statementBegin = true
}

// reportError reports an error to the diagnostics bag, starting at the
// passed position and ending at the parsers current position.
func (parsing *Parsing) reportError(error *diagnostic.RichError, region input.Region) {
	parsing.recorder.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SyntacticalAnalysis,
		UnitName: parsing.unitName,
		Position: region,
		Error:    error,
	})
}

type parsingError struct {
	Structure tree.NodeKind
	Position  input.Region
	Cause     error
}

func (err *parsingError) String() string {
	return fmt.Sprintf("failed parsing %s: %s",
		err.Structure.Name(),
		err.Cause.Error())
}

func (parsing *Parsing) throwError(cause *diagnostic.RichError) {
	structure, err := parsing.structureStack.pop()
	if err != nil {
		structure = structureStackElement{nodeKind: tree.UnknownNodeKind}
		log.Print("Could not pop structure stack")
	}
	region := input.CreateRegion(structure.beginOffset, parsing.offset())
	parsing.reportError(cause, region)
	panic(&parsingError{
		Structure: structure.nodeKind,
		Position:  region,
		Cause:     fmt.Errorf("could not parse %s", structure.nodeKind),
	})
}
func newNoIdentifierError(token token.Token) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "Identifier",
			Received: token.Value(),
		},
		CommonReasons: []string{
			"Declarations are not written properly",
		},
	}
}

func newInvalidKeywordError(token token.Token, expected token.Keyword) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: expected.String(),
			Received: token.Value(),
		},
		CommonReasons: aggregateReasonsOfInvalidKeywordError(token, expected),
	}
}

const forgottenDoKeywordReason = "The 'do' keyword after control structures is forgotten"

func aggregateReasonsOfInvalidKeywordError(
	received token.Token,
	expected token.Keyword) (reasons []string) {

	if expected == token.DoKeyword && token.IsEndOfStatementToken(received) {
		reasons = append(reasons, forgottenDoKeywordReason)
	}
	return reasons
}

func newInvalidOperatorError(token token.Token, expected token.Operator) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: expected.String(),
			Received: token.Value(),
		},
		CommonReasons: aggregateReasonsOfInvalidOperatorError(token),
	}
}

const unfinishedOperationReason = "An operation has not been completed"
const invalidOperator = "An invalid operator is applied to an operation"

func aggregateReasonsOfInvalidOperatorError(received token.Token) (reasons []string) {

	if token.IsIdentifierToken(received) {
		reasons = append(reasons, unfinishedOperationReason)
	}
	reasons = append(reasons, invalidOperator)
	return reasons
}

func newSmallerIndentError(indent token.Indent) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.InvalidIndentationError{
			Expected: "increased indent",
			Received: int(indent),
		},
	}
}
