package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type CallExpression struct {
	// Target is the procedure that is called.
	Target Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments    []*CallArgument
	Region input.Region
}

func (call *CallExpression) Accept(visitor Visitor) {
	VisitCallExpression(call)
}

func (call *CallExpression) AcceptRecursive(visitor Visitor) {
	call.Accept(visitor)
	AcceptRecursive(visitor)
	for _, argument := range call.Arguments {
		argument.AcceptRecursive(visitor)
	}
}

func (call *CallExpression) Locate() input.Region {
	return call.Region
}
