package diagnostic

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type RichError struct {
	Error         KnownError
	CommonReasons []string
}

type KnownError interface {
	Name() string
}

type NameCollisionError struct {
	Symbol string
}

func (error *NameCollisionError) Name() string {
	return fmt.Sprintf("name %s collides with another entry in the same scope", error.Symbol)
}

type UnexpectedTokenError struct {
	Expected string
	Received string
}

func (error *UnexpectedTokenError) Name() string {
	if error.Received != "" {
		return fmt.Sprintf("expected %s but got nothing", error.Expected)
	}
	return fmt.Sprintf("expected %s but got %s", error.Expected, error.Received)
}

type InvalidStatementError struct {
	Kind tree.NodeKind
}

func (error *InvalidStatementError) Name() string {
	return fmt.Sprintf("invalid statement of kind %s", error.Kind.Name())
}

type InvalidIndentationError struct {
	Expected string
	Received int
}

func (error *InvalidIndentationError) Name() string {
	if len(error.Expected) == 0 {
		return fmt.Sprintf("expected no indent but has %d", error.Received)
	}
	return fmt.Sprintf(
		"expected indent of %d but got %d", len(error.Expected), error.Received)
}

type SpecificError struct {
	Message string
}

func (error *SpecificError) Name() string {
	return error.Message
}
