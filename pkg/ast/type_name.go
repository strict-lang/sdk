package ast

import "fmt"

type TypeName interface {
	FullName() string
}

type ConcreteTypeName struct {
	Name string
}

func (concrete ConcreteTypeName) FullName() string {
	return concrete.Name
}

type GenericTypeName struct {
	Name    string
	Generic TypeName
}

func (generic GenericTypeName) FullName() string {
	return fmt.Sprintf("%s<%s>", generic.Name, generic.Generic)
}