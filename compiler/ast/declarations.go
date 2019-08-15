package ast

type ParameterList []*Parameter
type MethodDeclaration struct {
	Name       Identifier
	Type       TypeName
	Parameters ParameterList
	Body       Node
	NodePosition Position
}

func (method *MethodDeclaration) Accept(visitor *Visitor) {
	visitor.VisitMethod(method)
}

func (method *MethodDeclaration) AcceptAll(visitor *Visitor) {
	visitor.VisitMethod(method)
	for _, parameter := range method.Parameters {
		parameter.AcceptAll(visitor)
	}
	method.Body.AcceptAll(visitor)
}

func (method *MethodDeclaration) Position() Position {
	return method.Position()
}

type Parameter struct {
	Type TypeName
	Name Identifier
	NodePosition Position
}

func (parameter Parameter) IsNamedAfterType() bool {
	return parameter.Type.NonGenericName() == parameter.Name.Value
}

func (parameter *Parameter) Accept(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptAll(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) Position() Position {
	return parameter.NodePosition
}