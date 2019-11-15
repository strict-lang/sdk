package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type CallExpression struct {
	// Target is the procedure that is called.
	Target Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments    CallArgumentList
	Region       input.Region
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

func (call *CallExpression) Matches(node Node) bool {
	if target, ok := node.(*CallExpression); ok {
		return call.Target.Matches(target.Target) &&
			call.hasArguments(target.Arguments)
	}
	return false
}

func (call *CallExpression) hasArguments(arguments CallArgumentList) bool {
	if len(call.Arguments) != len(arguments) {
		return false
	}
	for index := 0; index < len(arguments); index++ {
		if !call.Arguments[index].Matches(arguments[index]) {
			return false
		}
	}
	return true
}
