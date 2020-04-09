package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

type ConstructorDeclaration struct {
	Parameters ParameterList
	Body       *StatementBlock
	Region     input.Region
	scope      scope.Scope
	Parent     Node
}

func (declaration *ConstructorDeclaration) UpdateScope(target scope.Scope) {
	declaration.scope = target
}

func (declaration *ConstructorDeclaration) Scope() scope.Scope {
	return declaration.scope
}

func (declaration *ConstructorDeclaration) SetEnclosingNode(target Node) {
	declaration.Parent = target
}

func (declaration *ConstructorDeclaration) EnclosingNode() (Node, bool) {
	return declaration.Parent, declaration.Parent != nil
}

func (declaration *ConstructorDeclaration) Accept(visitor Visitor) {
	visitor.VisitConstructorDeclaration(declaration)
}

func (declaration *ConstructorDeclaration) AcceptRecursive(visitor Visitor) {
	declaration.Accept(visitor)
	for _, parameter := range declaration.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	declaration.Body.AcceptRecursive(visitor)
}

func (declaration *ConstructorDeclaration) Locate() input.Region {
	return declaration.Region
}

func (declaration *ConstructorDeclaration) Matches(node Node) bool {
	if target, ok := node.(*ConstructorDeclaration); ok {
		return declaration.Parameters.Matches(target.Parameters) &&
			declaration.Body.Matches(target.Body)
	}
	return false
}
