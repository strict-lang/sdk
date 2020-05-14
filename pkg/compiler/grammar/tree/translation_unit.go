package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

// TranslationUnit represents a unit of translation, a file containing Strict
// input code. It can have multiple children, which are seen as the roots
// of the actual tree. This node however, is the real unit of the tree.
type TranslationUnit struct {
	Name      string
	Namespace string
	Imports   []*ImportStatement
	Class     *ClassDeclaration
	LineMap   *linemap.LineMap
	Region    input.Region
	scope     scope.Scope
	Trait     bool
}

func (unit *TranslationUnit) UpdateScope(target scope.Scope) {
	unit.scope = target
}

func (unit *TranslationUnit) Scope() scope.Scope {
	return unit.scope
}

func (*TranslationUnit) SetEnclosingNode(Node) {}

func (unit *TranslationUnit) EnclosingNode() (Node, bool) {
	return nil, false
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

func (unit *TranslationUnit) Matches(node Node) bool {
	if target, ok := node.(*TranslationUnit); ok {
		return unit.Name == target.Name &&
			unit.MatchesImports(target.Imports) &&
			unit.Class.Matches(target.Class)
	}
	return false
}

func (unit *TranslationUnit) MatchesImports(imports []*ImportStatement) bool {
	if len(unit.Imports) != len(imports) {
		return false
	}
	for index := 0; index < len(imports); index++ {
		if !unit.Imports[index].Matches(imports[index]) {
			return false
		}
	}
	return true
}
