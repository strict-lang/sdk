package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

type ParameterList []*Parameter

type MethodDeclaration struct {
	Name         *Identifier
	Type         TypeName
	Parameters   ParameterList
	Body         Node
	NodePosition InputRegion
}

func (method *MethodDeclaration) Accept(visitor Visitor) {
	visitor.VisitMethodDeclaration(method)
}

func (method *MethodDeclaration) AcceptRecursive(visitor Visitor) {
	visitor.VisitMethodDeclaration(method)
	for _, parameter := range method.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	AcceptRecursive(visitor)
}

func (method *MethodDeclaration) Area() InputRegion {
	return method.Area()
}

type Parameter struct {
	Type         TypeName
	Name         *Identifier
	Region  input.Region
}

func (parameter *Parameter) Accept(visitor Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) Locate() input.Region {
	return parameter.Region
}

