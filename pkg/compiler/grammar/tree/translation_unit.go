package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// TranslationUnit represents a unit of translation, a file containing Strict
// input code. It can have multiple children, which are seen as the roots
// of the actual tree. This node however, is the real unit of the tree.
type TranslationUnit struct {
	Name       string
	Imports    []*ImportStatement
	Class      *ClassDeclaration
	Region input.Region
}

func (unit *TranslationUnit) Accept(visitor Visitor) {
	visitor.VisitTranslationUnit(unit)
}

func (unit *TranslationUnit) AcceptRecursive(visitor Visitor) {
	unit.Accept(visitor)
	for _, importStatement := range unit.Imports {
		importStatement.AcceptRecursive(visitor)
	}
	unit.Class.AcceptRecursive(visitor)
}

func (unit *TranslationUnit) ToTypeName() TypeName {
	return &ConcreteTypeName{
		Name:   unit.Name,
		Region: unit.Region,
	}
}

func (unit *TranslationUnit) Locate() input.Region {
	return unit.Region
}
