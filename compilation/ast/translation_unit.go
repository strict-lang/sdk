package ast

// TranslationUnit represents a unit of translation, a file containing Strict
// source code. It can have multiple children, which are seen as the roots
// of the actual ast. This node however, is the real unit of the ast.
type TranslationUnit struct {
	name         string
	Children     []Node
	NodePosition Position
}

func NewEmptyTranslationUnit(name string) *TranslationUnit {
	return NewTranslationUnit(name, []Node{})
}

func NewTranslationUnit(name string, children []Node) *TranslationUnit {
	return &TranslationUnit{
		name:     name,
		Children: children,
	}
}

func (unit *TranslationUnit) Name() string {
	return unit.name
}

func (unit *TranslationUnit) Accept(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
}

func (unit *TranslationUnit) AcceptAll(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
	for _, topLevelNode := range unit.Children {
		topLevelNode.AcceptAll(visitor)
	}
}

func (unit *TranslationUnit) AppendChild(node Node) {
	unit.Children = append(unit.Children, node)
}

func (unit *TranslationUnit) ToTypeName() TypeName {
	return &ConcreteTypeName{
		Name:         unit.name,
		NodePosition: unit.NodePosition,
	}
}

func (unit *TranslationUnit) Position() Position {
	return unit.NodePosition
}
