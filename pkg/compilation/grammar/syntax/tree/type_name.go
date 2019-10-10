package tree

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
)

type TypeName interface {
	Node
	FullName() string
	NonGenericName() string
}

type ConcreteTypeName struct {
	Name         string
	Region input.Region
}

func (concrete *ConcreteTypeName) NonGenericName() string {
	return concrete.Name
}

func (concrete *ConcreteTypeName) FullName() string {
	return concrete.Name
}

func (concrete *ConcreteTypeName) Accept(visitor Visitor) {
	visitor.VisitConcreteTypeName(concrete)
}

func (concrete *ConcreteTypeName) AcceptRecursive(visitor Visitor) {
	visitor.VisitConcreteTypeName(concrete)
}

func (concrete *ConcreteTypeName) Locate() input.Region {
	return concrete.NodePosition
}

type GenericTypeName struct {
	Name string
	// TODO(merlinosayimwen) Change this to a slice to support
	//  types like maps and tuples.
	Generic      TypeName
	NodePosition InputRegion
}

func (generic *GenericTypeName) FullName() string {
	return fmt.Sprintf("%s<%s>", generic.Name, generic.Generic.FullName())
}

func (generic *GenericTypeName) NonGenericName() string {
	return generic.Name
}

func (generic *GenericTypeName) Accept(visitor *Visitor) {
	visitor.VisitGenericTypeName(generic)
}

func (generic *GenericTypeName) AcceptRecursive(visitor *Visitor) {
	visitor.VisitGenericTypeName(generic)
}

func (generic *GenericTypeName) Area() InputRegion {
	return generic.NodePosition
}

type ListTypeName struct {
	ElementTypeName TypeName
	NodePosition    InputRegion
}

func (list *ListTypeName) FullName() string {
	return fmt.Sprintf("%s[]", list.ElementTypeName.FullName())
}

func (list *ListTypeName) NonGenericName() string {
	return list.ElementTypeName.NonGenericName()
}

func (list *ListTypeName) Accept(visitor *Visitor) {
	visitor.VisitListTypeName(list)
}

func (list *ListTypeName) AcceptRecursive(visitor *Visitor) {
	visitor.VisitListTypeName(list)
	AcceptRecursive(visitor)
}

func (list *ListTypeName) Area() InputRegion {
	return list.NodePosition
}
