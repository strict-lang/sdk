package ast

import "fmt"

type TypeName interface {
	FullName() string
	NonGenericName() string
}

type ConcreteTypeName struct {
	Name string
}

func (concrete ConcreteTypeName) NonGenericName() string {
	return concrete.Name
}

func (concrete ConcreteTypeName) FullName() string {
	return concrete.Name
}

type GenericTypeName struct {
	Name    string
	Generic TypeName
}

func (generic GenericTypeName) FullName() string {
	return fmt.Sprintf("%s<%s>", generic.Name, generic.Generic.FullName())
}

func (generic GenericTypeName) NonGenericName() string {
	return generic.Name
}
