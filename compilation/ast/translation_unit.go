package ast

// TranslationUnit represents a unit of translation, a file containing Strict
// source code. It can have multiple children, which are seen as the roots
// of the actual ast. This node however, is the real unit of the ast.
type TranslationUnit struct {
	Name         string
	Imports      []*ImportStatement
	Class        ClassDeclaration
	NodePosition Position
}

func (unit *TranslationUnit) Accept(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
}

func (unit *TranslationUnit) AcceptRecursive(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
	for _, importStatement := range unit.Imports {
		importStatement.AcceptRecursive(visitor)
	}
	unit.Class.AcceptRecursive(visitor)
}

func (unit *TranslationUnit) ToTypeName() TypeName {
	return &ConcreteTypeName{
		Name:         unit.Name,
		NodePosition: unit.NodePosition,
	}
}

func (unit *TranslationUnit) Position() Position {
	return unit.NodePosition
}
