package syntaxtree

type ParameterList []*Parameter

type MethodDeclaration struct {
	Name         *Identifier
	Type         TypeName
	Parameters   ParameterList
	Body         Node
	NodePosition Position
}

func (method *MethodDeclaration) Accept(visitor *Visitor) {
	visitor.VisitMethodDeclaration(method)
}

func (method *MethodDeclaration) AcceptRecursive(visitor *Visitor) {
	visitor.VisitMethodDeclaration(method)
	for _, parameter := range method.Parameters {
		parameter.AcceptRecursive(visitor)
	}
	method.Body.AcceptRecursive(visitor)
}

func (method *MethodDeclaration) Position() Position {
	return method.Position()
}

type Parameter struct {
	Type         TypeName
	Name         *Identifier
	NodePosition Position
}

func (parameter Parameter) IsNamedAfterType() bool {
	return parameter.Type.NonGenericName() == parameter.Name.Value
}

func (parameter *Parameter) Accept(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) Position() Position {
	return parameter.NodePosition
}

type FieldDeclaration struct {
	Name         *Identifier
	TypeName     TypeName
	NodePosition Position
}

func (field *FieldDeclaration) Accept(visitor *Visitor) {
	visitor.VisitFieldDeclaration(field)
}

func (field *FieldDeclaration) AcceptRecursive(visitor *Visitor) {
	visitor.VisitFieldDeclaration(field)
}

func (field *FieldDeclaration) Position() Position {
	return field.NodePosition
}

type ClassDeclaration struct {
	Name         string
	Parameters   []ClassParameter
	SuperTypes   []TypeName
	Children     []Node
	NodePosition Position
}

func (class *ClassDeclaration) Accept(visitor *Visitor) {
	visitor.VisitClassDeclaration(class)
}

func (class *ClassDeclaration) AcceptRecursive(visitor *Visitor) {
	visitor.VisitClassDeclaration(class)
	for _, child := range class.Children {
		child.AcceptRecursive(visitor)
	}
}

func (class *ClassDeclaration) Position() Position {
	return class.NodePosition
}

type ClassParameter struct {
	Name      string
	SuperType TypeName
}

type ConstructorDeclaration struct {
	Parameters ParameterList
	Body Node
	NodePosition Position
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

func (declaration *ConstructorDeclaration) Position() Position {
	return declaration.NodePosition
}