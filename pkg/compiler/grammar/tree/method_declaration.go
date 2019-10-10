package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ParameterList []*Parameter

type MethodDeclaration struct {
	Name       *Identifier
	Type       TypeName
	Parameters ParameterList
	Body       Node
	Region     input.Region
}

func (declaration *MethodDeclaration) Accept(visitor Visitor) {
	VisitMethodDeclaration(declaration)
}

func (declaration *MethodDeclaration) AcceptRecursive(visitor Visitor) {
	VisitMethodDeclaration(declaration)
	for _, parameter := range declaration.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	declaration.AcceptRecursive(visitor)
}

func (declaration *MethodDeclaration) Locate() input.Region {
	return declaration.Region
}