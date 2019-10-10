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

func (method *MethodDeclaration) Accept(visitor *Visitor) {
	visitor.VisitMethodDeclaration(method)
}

func (method *MethodDeclaration) AcceptRecursive(visitor *Visitor) {
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
	NodePosition InputRegion
}

func (parameter *Parameter) Accept(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) Area() InputRegion {
	return parameter.NodePosition
}

type FieldDeclaration struct {
	Name         *Identifier
	TypeName     TypeName
	NodePosition InputRegion
}

func (field *FieldDeclaration) Accept(visitor *Visitor) {
	visitor.VisitFieldDeclaration(field)
}

func (field *FieldDeclaration) AcceptRecursive(visitor *Visitor) {
	visitor.VisitFieldDeclaration(field)
}

func (field *FieldDeclaration) Area() InputRegion {
	return field.NodePosition
}

type ClassDeclaration struct {
	Name         string
	Parameters   []ClassParameter
	SuperTypes   []TypeName
	Children     []Node
	NodeRegion   input.Region
}

func (class *ClassDeclaration) Accept(visitor *Visitor) {
	visitor.VisitClassDeclaration(class)
}

func (class *ClassDeclaration) AcceptRecursive(visitor *Visitor) {
	visitor.VisitClassDeclaration(class)
	for _, child := range class.Children {
		AcceptRecursive(visitor)
	}
}

func (class *ClassDeclaration) Area() InputRegion {
	return class.NodePosition
}

type ClassParameter struct {
	Name      string
	SuperType TypeName
}

type ConstructorDeclaration struct {
	Parameters   ParameterList
	Body         Node
	NodePosition InputRegion
}

func (declaration *ConstructorDeclaration) Accept(visitor *Visitor) {
	visitor.VisitConstructorDeclaration(declaration)
}

func (declaration *ConstructorDeclaration) AcceptRecursive(visitor *Visitor) {
	visitor.VisitConstructorDeclaration(declaration)
	for _, parameter := range declaration.Parameters {
		parameter.AcceptRecursive(visitor)
	}
}

func (declaration *ConstructorDeclaration) Area() InputRegion {
	return declaration.NodePosition
}
