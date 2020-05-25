package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

type ParameterList []*Parameter

type MethodDeclaration struct {
	Name       *Identifier
	Type       TypeName
	Parameters ParameterList
	Body       Node
	Region     input.Region
	Parent     Node
	Abstract   bool
	Factory    bool
	scope      scope.Scope
}

func (declaration *MethodDeclaration) UpdateScope(target scope.Scope) {
	declaration.scope = target
}

func (declaration *MethodDeclaration) Scope() scope.Scope {
	return declaration.scope
}

func (declaration *MethodDeclaration) SetEnclosingNode(target Node) {
	declaration.Parent = target
}

func (declaration *MethodDeclaration) EnclosingNode() (Node, bool) {
	return declaration.Parent, declaration.Parent != nil
}

func (declaration *MethodDeclaration) Accept(visitor Visitor) {
	visitor.VisitMethodDeclaration(declaration)
}

func (declaration *MethodDeclaration) AcceptRecursive(visitor Visitor) {
	declaration.Accept(visitor)
	if declaration.Type != nil {
		declaration.Type.AcceptRecursive(visitor)
	}
	for _, parameter := range declaration.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	declaration.Body.AcceptRecursive(visitor)
}

func (declaration *MethodDeclaration) Locate() input.Region {
	return declaration.Region
}

func (declaration *MethodDeclaration) Matches(node Node) bool {
	if target, ok := node.(*MethodDeclaration); ok {
		return declaration.Factory == target.Factory &&
			declaration.Name.Matches(target.Name) &&
			declaration.Type.Matches(target.Type) &&
			declaration.Parameters.Matches(target.Parameters) &&
			declaration.Body.Matches(target.Body)
	}
	return false
}

func (list ParameterList) Matches(target ParameterList) bool {
	if len(list) != len(target) {
		return false
	}
	for index := 0; index < len(list); index++ {
		if !list[index].Matches(target[index]) {
			return false
		}
	}
	return true
}
