package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ConstructorDeclaration struct {
	Parameters ParameterList
	Child      Node
	Region     input.Region
}

func (declaration *ConstructorDeclaration) Accept(visitor Visitor) {
	visitor.VisitConstructorDeclaration(declaration)
}

func (declaration *ConstructorDeclaration) AcceptRecursive(visitor Visitor) {
	declaration.Accept(visitor)
	for _, parameter := range declaration.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	declaration.Child.AcceptRecursive(visitor)
}

func (declaration *ConstructorDeclaration) Locate() input.Region {
	return declaration.Region
}


func (declaration *ConstructorDeclaration) Matches(node Node) bool {
	if target, ok := node.(*ConstructorDeclaration); ok {
		return declaration.Parameters.Matches(target.Parameters) &&
			declaration.Child.Matches(target.Child)
	}
	return false
}