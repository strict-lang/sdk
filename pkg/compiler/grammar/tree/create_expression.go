package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/scope"
)

type CreateExpression struct {
	Region       input.Region
	Call         *CallExpression
	Type         TypeName
	resolvedType resolvedType
	Parent Node
}

func (create *CreateExpression) SetEnclosingNode(target Node) {
  create.Parent = target
}

func (create *CreateExpression) EnclosingNode() (Node, bool) {
  return create.Parent, create.Parent != nil
}

func (create *CreateExpression) ResolveType(class *scope.Class) {
  create.resolvedType.resolve(class)
}

func (create *CreateExpression) ResolvedType() (*scope.Class, bool) {
  return create.resolvedType.class()
}

func (create *CreateExpression) Accept(visitor Visitor) {
	visitor.VisitCreateExpression(create)
}

func (create *CreateExpression) AcceptRecursive(visitor Visitor) {
	create.Accept(visitor)
	create.Type.AcceptRecursive(visitor)
	create.Call.AcceptRecursive(visitor)
}

func (create *CreateExpression) Locate() input.Region {
	return create.Region
}

func (create *CreateExpression) Matches(node Node) bool {
	if target, ok := node.(*CreateExpression); ok {
		return create.Call.Matches(target.Call) &&
			create.Type.Matches(target.Type)
	}
	return false
}

func (create *CreateExpression) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteCreateExpression(create)
}