package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
	"gitlab.com/strict-lang/sdk/pkg/compiler/typing"
)

type ClassDeclaration struct {
	Name       string
	Parameters []*ClassParameter
	SuperTypes []TypeName
	Children   []Node
	Region     input.Region
	Parent     Node
	scope      scope.Scope
}

type ClassParameter struct {
	Name      string
	SuperType TypeName
	Parent    Node
}

func (class *ClassDeclaration) UpdateScope(target scope.Scope) {
	class.scope = target
}

func (class *ClassDeclaration) Scope() scope.Scope {
	return class.scope
}

func (parameter *ClassParameter) SetEnclosingNode(target Node) {
	parameter.Parent = target
}

func (parameter *ClassParameter) EnclosingNode() (Node, bool) {
	return parameter.Parent, parameter.Parent != nil
}

func (parameter *ClassParameter) Matches(target *ClassParameter) bool {
	return parameter.Name == target.Name &&
		parameter.SuperType.Matches(target.SuperType)
}

func (class *ClassDeclaration) SetEnclosingNode(target Node) {
	class.Parent = target
}

func (class *ClassDeclaration) EnclosingNode() (Node, bool) {
	return class.Parent, class.Parent != nil
}

func (class *ClassDeclaration) Accept(visitor Visitor) {
	visitor.VisitClassDeclaration(class)
}

func (class *ClassDeclaration) AcceptRecursive(visitor Visitor) {
	class.Accept(visitor)
	for _, superType := range class.SuperTypes {
		superType.AcceptRecursive(visitor)
	}
	for _, child := range class.Children {
		child.AcceptRecursive(visitor)
	}
}

func (class *ClassDeclaration) Locate() input.Region {
	return class.Region
}

func (class *ClassDeclaration) Matches(node Node) bool {
	if target, ok := node.(*ClassDeclaration); ok {
		return class.matchesClass(target)
	}
	return false
}

func (class *ClassDeclaration) matchesClass(target *ClassDeclaration) bool {
	return class.Name == target.Name &&
		class.hasParameters(target.Parameters) &&
		class.hasSuperTypes(target.SuperTypes) &&
		class.hasChildren(target.Children)
}

func (class *ClassDeclaration) hasParameters(parameters []*ClassParameter) bool {
	if len(class.Parameters) != len(parameters) {
		return false
	}
	for index := 0; index < len(parameters); index++ {
		if class.Parameters[index].Matches(parameters[index]) {
			return false
		}
	}
	return true
}

func (class *ClassDeclaration) hasSuperTypes(types []TypeName) bool {
	if len(class.SuperTypes) != len(types) {
		return false
	}
	for index := 0; index < len(types); index++ {
		if class.SuperTypes[index].Matches(types[index]) {
			return false
		}
	}
	return true
}

func (class *ClassDeclaration) hasChildren(children []Node) bool {
	if len(class.Children) != len(children) {
		return false
	}
	for index := 0; index < len(children); index++ {
		if !class.Children[index].Matches(children[index]) {
			return false
		}
	}
	return true
}

func (class *ClassDeclaration) NewActualClass() typing.Type {
	// TODO: Create proper class
	return &typing.ConcreteType{
		Name:   class.Name,
		Traits: []typing.Type{},
	}
}
