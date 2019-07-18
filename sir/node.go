package sir

import (
	"gitlab.com/strict-lang/sdk/sir/constantpool"
	"gitlab.com/strict-lang/sdk/sir/metadata"
)

type Node interface {
	Accept(visitor *Visitor)
	AcceptAll(visitor *Visitor)
}

type Unit struct {
	Module 			 *Module
	ConstantPool *constantpool.Pool
}

func (unit *Unit) Accept(visitor *Visitor) {
	visitor.VisitUnit(unit)
}

func (unit *Unit) AcceptAll(visitor *Visitor) {
	visitor.VisitUnit(unit)
	unit.Module.AcceptAll(visitor)
}

type Typed struct {
	TypeName constantpool.Reference
}

type Declaration struct {
	Node
	Name     constantpool.Reference
	Metadata *metadata.Table
}

type TypedDeclaration struct {
	Typed
	Declaration
}

type Module struct {
	Declaration
	TopLevelDeclarations []Declaration
}

func (module *Module) Accept(visitor *Visitor) {
	visitor.VisitModule(module)
}

func (module *Module) AcceptAll(visitor *Visitor) {
	visitor.VisitModule(module)
	for _, declaration := range module.TopLevelDeclarations {
		declaration.AcceptAll(visitor)
	}
}

type MethodDeclaration struct {
	TypedDeclaration
	Parameters []*MethodParameter
	CodeBlock *CodeBlock
}

func (method *MethodDeclaration) Accept(visitor *Visitor) {
	visitor.VisitMethodDeclaration(method)
}

func (method *MethodDeclaration) AcceptAll(visitor *Visitor) {
	visitor.VisitMethodDeclaration(method)
	for _, parameter := range method.Parameters {
		parameter.AcceptAll(visitor)
	}
}

type MethodParameter struct {
	TypedDeclaration
}

func (parameter *MethodParameter) Accept(visitor *Visitor) {
	visitor.VisitMethodParameter(parameter)
}

func (parameter *MethodParameter) AcceptAll(visitor *Visitor) {
	visitor.VisitMethodParameter(parameter)
}
