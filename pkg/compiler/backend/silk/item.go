package silk

import (
	"errors"
	"strict.dev/sdk/pkg/compiler/typing"
	"strict.dev/sdk/pkg/silk"
)

type Item struct {
	Location StorageLocation
	Name string
	Type typing.Type
	Context *MethodContext
}

type StorageLocation interface {
	EmitLoad(item *Item, generation *Generation) error
	EmitStore(item *Item, generation *Generation) error
	CreateSilkLocation() silk.StorageLocation
}

type VariableLocation struct {
	Index int
	Name string
	Type typing.Type
}

func (location *VariableLocation) EmitLoad(item *Item, generation *Generation) error {
	if !location.canEmitLoad(item) {
		return errors.New("can not load item from target location")
	}
	target, err := generation.selectLoadTarget(item)
	if err != nil {
		return err
	}
	return location.loadIntoTarget(item, generation, target)
}

func (location *VariableLocation) CreateSilkLocation() silk.StorageLocation {
	return &silk.Field{
		Name: location.Name,
		Index: location.Index,
		Type: translateType(location.Type),
	}
}

func (location *VariableLocation) loadIntoTarget(
	item *Item, generation *Generation, target StorageLocation) error {

	item.Location = location
	generation.EmitInstruction(silk.Instruct(&silk.LoadOperation{
		Type: translateType(location.Type),
		Index: location.Index,
		Target: target.CreateSilkLocation(),
	}))
	return nil
}

func (location *VariableLocation) canEmitLoad(item *Item) bool {
	variable := item.Context.IndexVariable(location.Index)
	return isAssignCompatible(variable.Type, item.Type)
}

func (location *VariableLocation) EmitStore(item *Item, generation *Generation) error {
	if !isAssignCompatible(location.Type, item.Type) {
		return errors.New("can not store item to target location")
	}
	generation.EmitInstruction(silk.Instruct(&silk.StoreOperation{
		Type: translateType(location.Type),
		Index: location.Index,
		Target: location.CreateSilkLocation(),
	}))
	return nil
}

func isAssignCompatible(left typing.Type, right typing.Type) bool {
	return left.Is(right)
}

type StackLocation struct {
	Offset int
}
