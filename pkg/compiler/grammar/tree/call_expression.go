package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type CallExpression struct {
	// Target is the procedure that is called.
	Target Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments CallArgumentList
	Region    input.Region
	resolvedType resolvedType
}

func (call *CallExpression) Resolve(descriptor TypeDescriptor) {
	call.resolvedType.setDescriptor(descriptor)
}

func (call *CallExpression) GetResolvedType() (TypeDescriptor, bool) {
	return call.resolvedType.getDescriptor()
}

func (call *CallExpression) Accept(visitor Visitor) {
	visitor.VisitCallExpression(call)
}

func (call *CallExpression) AcceptRecursive(visitor Visitor) {
	call.Accept(visitor)
	call.Target.AcceptRecursive(visitor)
	for _, argument := range call.Arguments {
		argument.AcceptRecursive(visitor)
	}
}

func (call *CallExpression) Locate() input.Region {
	return call.Region
}
