package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

type CallExpression struct {
	// Target is the procedure that is called.
	Target Expression
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments    CallArgumentList
	Region       input.Region
	resolvedType resolvedType
	Parent       Node
	name         cachedName
}

type cachedName struct {
	value    *Identifier
	found    bool
	searched bool
}

func (call *CallExpression) SetEnclosingNode(target Node) {
	call.Parent = target
}

func (call *CallExpression) EnclosingNode() (Node, bool) {
	return call.Parent, call.Parent != nil
}

func (call *CallExpression) ResolveType(class *scope.Class) {
	call.resolvedType.resolve(class)
}

func (call *CallExpression) ResolvedType() (*scope.Class, bool) {
	return call.resolvedType.class()
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

func (call *CallExpression) TargetName() *Identifier {
	call.maybeTryToResolveName()
	return call.name.value
}

func (call *CallExpression) IsNamedCall() bool {
	call.maybeTryToResolveName()
	return call.name.found
}

func (call *CallExpression) maybeTryToResolveName() {
	if !call.name.searched {
		call.resolveName()
		call.name.searched = true
	}
}

func (call *CallExpression) resolveName() {
	visitor := call.createNameResolveVisitor()
	call.Target.Accept(visitor)
}

func (call *CallExpression) createNameResolveVisitor() Visitor {
	return &DelegatingVisitor{
		IdentifierVisitor: func(identifier *Identifier) {
			call.name.found = true
			call.name.value = identifier
		},
		FieldSelectExpressionVisitor: func(expression *FieldSelectExpression) {
			identifier, found := expression.FindLastIdentifier()
			call.name.found = found
			call.name.value = identifier
		},
	}
}

func (call *CallExpression) TransformExpressions(transformer ExpressionTransformer) {
	call.Target = call.Target.Transform(transformer)
}

func (call *CallExpression) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteCallExpression(call)
}
