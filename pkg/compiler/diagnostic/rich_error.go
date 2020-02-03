package diagnostic

import "strict.dev/sdk/pkg/compiler/grammar/tree"

type RichError struct {
	Error         KnownError
	CommonReasons []string
}

type KnownError interface {
	Accept(visitor ErrorVisitor)
}

type ErrorVisitor interface {
	VisitUnexpectedToken(*UnexpectedTokenError)
	VisitInvalidStatement(*InvalidStatementError)
	VisitInvalidIndentation(*InvalidIndentationError)
	VisitSpecificError(*SpecificError)
}

type UnexpectedTokenError struct {
	Expected string
	Received string
}

func (error *UnexpectedTokenError) Accept(visitor ErrorVisitor) {
	visitor.VisitUnexpectedToken(error)
}

type InvalidStatementError struct {
	Kind tree.NodeKind
}

func (error *InvalidStatementError) Accept(visitor ErrorVisitor) {
	visitor.VisitInvalidStatement(error)
}

type InvalidIndentationError struct {
	Expected string
	Received int
}

func (error *InvalidIndentationError) Accept(visitor ErrorVisitor) {
	visitor.VisitInvalidIndentation(error)
}

type SpecificError struct {
	Message string
}

func (error *SpecificError) Accept(visitor ErrorVisitor) {
	visitor.VisitSpecificError(error)
}
