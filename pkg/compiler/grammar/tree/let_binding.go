package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

// LetBinding binds an expression to a name. In strict, this behaves like a
// variable definition. LetBindings are also viable expressions, which have the
// value of their Expression field.
type LetBinding struct {
	Parent Node
	Region input.Region
	Expression Expression
	Name *Identifier
}

func (binding *LetBinding) ResolveType(class *scope.Class) {
	binding.Expression.ResolveType(class)
	binding.Name.ResolveType(class)
}

func (binding *LetBinding) ResolvedType() (*scope.Class, bool) {
	return binding.Expression.ResolvedType()
}

func (binding *LetBinding) Locate() input.Region {
	return binding.Region
}

func (binding *LetBinding) Accept(visitor Visitor) {}

func (binding *LetBinding) AcceptRecursive(visitor Visitor) {
	binding.Expression.AcceptRecursive(visitor)
	binding.Name.AcceptRecursive(visitor)
}

func (binding *LetBinding) SetEnclosingNode(target Node) {
	binding.Parent = target
}

func (binding *LetBinding) EnclosingNode() (Node, bool) {
	return binding.Parent, binding.Parent != nil
}

func (binding *LetBinding) Matches(target Node) bool {
	if _, isWildcard := target.(*WildcardNode); isWildcard {
		return true
	}
	targetBinding, ok := target.(*LetBinding)
	return ok && binding.matchesBinding(targetBinding)
}

func (binding *LetBinding) matchesBinding(target *LetBinding) bool {
	return binding.Name.Matches(target.Name)  &&
		binding.Expression.Matches(target.Expression)
}

func (binding *LetBinding) TransformExpressions(transformer ExpressionTransformer) {
	binding.Expression = binding.Expression.Transform(transformer)
}

func (binding *LetBinding) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteLetBinding(binding)
}