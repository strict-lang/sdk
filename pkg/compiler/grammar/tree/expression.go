package tree

import "strict.dev/sdk/pkg/compiler/scope"

type Expression interface {
	Node
	ResolveType(symbol *scope.Class)
	ResolvedType() (*scope.Class, bool)
	Transform(ExpressionTransformer) Expression
}

type StoredExpression interface {
	Expression
}

type ExpressionContainer interface {
	TransformExpressions(transformer ExpressionTransformer)
}
