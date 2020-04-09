package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

// LetBinding binds an expression to a name. In strict, this behaves like a
// variable definition. LetBindings are also viable expressions, which have the
// value of their Expression field.
type LetBinding struct {
	Parent     Node
	Region     input.Region
	Expression Expression
	Names[]       *Identifier
}

func (binding *LetBinding) ResolveType(class *scope.Class) {
	binding.Expression.ResolveType(class)
	for _, name := range binding.Names {
		name.ResolveType(class)
	}
}

func (binding *LetBinding) ResolvedType() (*scope.Class, bool) {
	return binding.Expression.ResolvedType()
}

func (binding *LetBinding) Locate() input.Region {
	return binding.Region
}

func (binding *LetBinding) Accept(visitor Visitor) {
	visitor.VisitLetBinding(binding)
}

func (binding *LetBinding) AcceptRecursive(visitor Visitor) {
	binding.Accept(visitor)
	binding.Expression.AcceptRecursive(visitor)
	for _, name := range binding.Names {
		name.AcceptRecursive(visitor)
	}
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
	return binding.matchesNames(target.Names) &&
		binding.Expression.Matches(target.Expression)
}

func (binding *LetBinding) matchesNames(names []*Identifier) bool {
	if len(names) != len(binding.Names) {
		return false
	}
	for index, name := range binding.Names {
		if names[index] != name {
			return false
		}
	}
	return true
}

func (binding *LetBinding) TransformExpressions(transformer ExpressionTransformer) {
	binding.Expression = binding.Expression.Transform(transformer)
}

func (binding *LetBinding) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteLetBinding(binding)
}
