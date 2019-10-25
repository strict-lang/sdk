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
	visitor.VisitMethodDeclaration(declaration)
}

func (declaration *MethodDeclaration) AcceptRecursive(visitor Visitor) {
	declaration.Accept(visitor)
	for _, parameter := range declaration.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	declaration.AcceptRecursive(visitor)
}

func (declaration *MethodDeclaration) Locate() input.Region {
	return declaration.Region
}

func (declaration *MethodDeclaration) Matches(node Node) bool {
	if target, ok := node.(*MethodDeclaration); ok {
		return declaration.Name.Matches(target.Name) &&
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