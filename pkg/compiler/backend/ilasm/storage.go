package ilasm

import (
	"errors"
)

type Item struct {
	Class *Class
}

type StorageLocation interface {
	EmitLoad(code *BlockBuilder) error
	EmitStore(item *Item, code *BlockBuilder) error
}

type StackLocation struct {
	Index int
}

func (location *StackLocation) EmitLoad(code *BlockBuilder) error {
	// Value is already on the stack.
	return nil
}

func (location *StackLocation) EmitStore(item *Item, code *BlockBuilder) error {
	code.EmitPush(item.Class)
	return nil
}

type MemberField struct {
	Name           string
	Class          *Class
	EnclosingClass *Class
}

type MemberLocation struct {
	Field MemberField
	InstanceLocation StorageLocation
}

func NewOwnMemberFieldLocation(ownClass *Class, field MemberField) StorageLocation {
	return &MemberLocation{
		Field:            field,
		InstanceLocation: createOwnReferenceLocation(ownClass),
	}
}

func createOwnReferenceLocation(ownClass *Class) StorageLocation {
	return &VariableLocation{
		Variable:  &VirtualVariable{
			Class: ownClass,
			Index: 0,
		},
		Parameter: true,
	}
}

func (location *MemberLocation) EmitLoad(code *BlockBuilder) error {
	if err := location.InstanceLocation.EmitLoad(code); err != nil {
		return err
	}
	code.EmitMemberLoad(location.Field)
	return nil
}

func (location *MemberLocation) EmitStore(item *Item, code *BlockBuilder) error {
	beforeLastPush := code.BeforeLastPush()
	if err := location.InstanceLocation.EmitLoad(beforeLastPush); err != nil {
		return err
	}
	code.EmitMemberStore(location.Field)
	return nil
}

type VariableLocation struct {
	Variable *VirtualVariable
	Parameter bool
}

func (location *VariableLocation) EmitLoad(code *BlockBuilder) error {
	code.EmitFieldLoad(location.Variable)
	return nil
}

func (location *VariableLocation) EmitStore(item *Item, code *BlockBuilder) error {
	targetClass := location.Variable.Class
	if !item.Class.IsAssignable(targetClass) {
		return newNotAssignableError(item, targetClass)
	}
	code.EmitFieldStore(location.Variable)
	return nil
}

func newNotAssignableError(item *Item, targetClass *Class) error {
	return errors.New("item is not assignable to variable")
}

type VirtualVariable struct {
	Class *Class
	Name string
	Index int
}

type VirtualOperandStack struct {
	CurrentDepth int
	MaximumDepth int
}

func (stack *VirtualOperandStack) IncreaseDepth() {
	updatedDepth := stack.CurrentDepth + 1
	if updatedDepth > stack.MaximumDepth {
		stack.MaximumDepth = updatedDepth
	}
	stack.CurrentDepth = updatedDepth
}

func (stack *VirtualOperandStack) DecreaseDepth() {
	if updatedDepth := stack.CurrentDepth - 1; updatedDepth >= 0 {
		stack.CurrentDepth = updatedDepth
	}
}

type MethodContext struct {
	Variables []*VariableSet
	Arguments []*VariableSet
	Stack *VirtualOperandStack
}

type VariableSet struct {
	indexed []*VirtualVariable
	named map[string] *VirtualVariable
}

func (set *VariableSet) Index(index int) *VirtualVariable {
	return set.indexed[index]
}

func (set *VariableSet) IndexOrCreate(index int, class *Class) *VirtualVariable {
	if len(set.indexed) > index {
		return set.indexOrCreateInBounds(index, class)
	}
	return set.createOutsideOfBounds(index, class)
}

func (set *VariableSet) indexOrCreateInBounds(index int, class *Class) *VirtualVariable {
	if variable := set.indexed[index]; variable != nil {
		return variable
	}
	created := &VirtualVariable{Index: index, Class: class}
	set.insert(index, created)
	return created
}

func (set *VariableSet) createOutsideOfBounds(index int, class *Class) *VirtualVariable {
	created := &VirtualVariable{Index: index, Class: class}
	set.expandAndInsert(index, created)
	return created
}

func (set *VariableSet) expandAndInsert(index int, variable *VirtualVariable) {
	newIndexed := make([]*VirtualVariable, index + 1)
	copy(newIndexed, set.indexed)
	newIndexed[index] = variable
	set.indexed = newIndexed
}

func (set *VariableSet) insert(index int, variable *VirtualVariable) {
	set.indexed[index] = variable
}

func (set *VariableSet) FindByName(name string) *VirtualVariable {
	return set.named[name]
}

func (set *VariableSet) FindByNameOrCreate(name string, class *Class) *VirtualVariable {
	if variable, ok := set.named[name]; ok {
		return variable
	}
	return set.createAndInsertNamed(name, class)
}

func (set *VariableSet) createAndInsertNamed(name string, class *Class) *VirtualVariable {
	created := &VirtualVariable{Name: name, Class: class}
	set.named[name] = created
	return created
}