package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/scope"

type Expression interface {
	Node
	ResolveType(symbol *scope.Class)
	ResolvedType() (*scope.Class, bool)
}

type StoredExpression interface {
	Expression
}