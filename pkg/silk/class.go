package silk

import "gitlab.com/strict-lang/sdk/pkg/silk/symbol"

type Class struct {
	Name    symbol.Reference
	Traits  []symbol.Reference
	Methods []*Method
}

type ClassBuilder struct {
	tableBuilder    *symbol.TableBuilder
	nameReference   symbol.Reference
	traitReferences []symbol.Reference
	methods         []*Method
}

func (builder *ClassBuilder) WithName(name string) *ClassBuilder {
	nameSymbol := &symbol.Name{Value: name}
	reference := builder.tableBuilder.Insert(nameSymbol)
	return builder.WithNameReference(reference)
}

func (builder *ClassBuilder) WithNameReference(
	reference symbol.Reference) *ClassBuilder {

	builder.nameReference = reference
	return builder
}

func (builder *ClassBuilder) AddTraitReference(
	reference symbol.Reference) *ClassBuilder {

	builder.traitReferences = append(builder.traitReferences, reference)
	return builder
}

func (builder *ClassBuilder) AddMethod(method *Method) *ClassBuilder {
	builder.methods = append(builder.methods, method)
	return builder
}

func (builder *ClassBuilder) Create() *Class {
	traits := make([]symbol.Reference, len(builder.traitReferences))
	copy(traits, builder.traitReferences)
	methods := make([]*Method, len(builder.methods))
	copy(methods, builder.methods)
	return &Class{
		Name:    builder.nameReference,
		Traits:  traits,
		Methods: methods,
	}
}
