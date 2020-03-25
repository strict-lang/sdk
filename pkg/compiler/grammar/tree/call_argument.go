package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

type CallArgumentList []*CallArgument

type CallArgument struct {
	Label  string
	Value  Expression
	Region input.Region
	Parent Node
}

func (argument *CallArgument) SetEnclosingNode(target Node) {
	argument.Parent = target
}

func (argument *CallArgument) EnclosingNode() (Node, bool) {
	return argument.Parent, argument.Parent != nil
}

func (argument *CallArgument) IsLabeled() bool {
	return argument.Label != ""
}

func (argument *CallArgument) Accept(visitor Visitor) {
	visitor.VisitCallArgument(argument)
}

func (argument *CallArgument) AcceptRecursive(visitor Visitor) {
	argument.Accept(visitor)
	argument.Value.AcceptRecursive(visitor)
}

func (argument *CallArgument) Locate() input.Region {
	return argument.Region
}

func (argument *CallArgument) Matches(node Node) bool {
	if target, ok := node.(*CallArgument); ok {
		return argument.Label == target.Label &&
			argument.Value.Matches(target.Value)
	}
	return false
}

func (argument *CallArgument) ResolvedType() (*scope.Class, bool) {
	return argument.Value.ResolvedType()
}

func (argument *CallArgument) ResolveType(class *scope.Class) {
	argument.Value.ResolveType(class)
}

func (argument *CallArgument) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteCallArgument(argument)
}

func (argument *CallArgument) TransformExpressions(transformer ExpressionTransformer) {
	argument.Value = argument.Value.Transform(transformer)
}
