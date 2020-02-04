package silk

import "strict.dev/sdk/pkg/silk/symbol"

type Operation interface {
	Matches(Operation) bool
	Accept(visitor Visitor, instruction *Instruction)
}

func isWildcard(operation Operation) bool {
	_, wildcard := operation.(WildcardOperation)
	return wildcard
}

type WildcardOperation struct{}

func (WildcardOperation) Matches(Operation) bool {
	return true
}

func (WildcardOperation) Accept(visitor Visitor, instruction *Instruction) {}

type NoOperation struct{}

func (operation NoOperation) Matches(target Operation) bool {
	_, isSameType := target.(NoOperation)
	return isSameType
}

func (operation NoOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitNoOperation(instruction, operation)
}

type Variable interface{}

type ArithmeticOperation int

type ArithmeticInstruction struct {
	Output    *Field
	Operands  []Variable
	Operation ArithmeticOperation
}

type StorageLocation interface {
	Matches(StorageLocation) bool
}

type LoadOperation struct {
	Index  int
	Type   Type
	Target StorageLocation
}

func (load *LoadOperation) Matches(operation Operation) bool {
	if isWildcard(operation) {
		return true
	}
	if targetLoad, ok := operation.(*LoadOperation); ok {
		return load.matchesDeeply(targetLoad)
	}
	return false
}

func (load *LoadOperation) matchesDeeply(target *LoadOperation) bool {
	return load.Index == target.Index &&
		load.Target.Matches(target.Target) &&
		load.Type.Matches(target.Type)

}

func (load *LoadOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitLoad(instruction, load)
}

type StoreOperation struct {
	Index  int
	Type   Type
	Target StorageLocation
}

func (store *StoreOperation) Matches(operation Operation) bool {
	if isWildcard(operation) {
		return true
	}
	if targetLoad, ok := operation.(*StoreOperation); ok {
		return store.matchesDeeply(targetLoad)
	}
	return false
}

func (store *StoreOperation) matchesDeeply(target *StoreOperation) bool {
	return store.Index == target.Index &&
		store.Target.Matches(target.Target) &&
		store.Type.Matches(target.Type)

}

func (store *StoreOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitStore(instruction, store)
}

type PushOperation struct {
	Type     Type
	Variable Variable
}

type PopOperation struct {
	Type Type
}

func (pop *PopOperation) Matches(operation Operation) bool {
	if isWildcard(operation) {
		return true
	}
	if targetPop, ok := operation.(*PopOperation); ok {
		return pop.matchesDeeply(targetPop)
	}
	return false
}

func (pop *PopOperation) matchesDeeply(target *PopOperation) bool {
	return pop.Type.Matches(target.Type)
}

func (pop *PopOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitPop(instruction, pop)
}

type Comparison int

type CompareInstruction struct {
	Output     *Field
	Comparison Comparison
	Operands   []Variable
}

type JumpOperation struct{}

type PhiEntry struct {
	LastBlock BlockId
	Value     Variable
}

type PhiOperation struct {
	ReturnedType Type
	Entries      []PhiEntry
	Block        *Block
	Output       *Field
}

func (operation *PhiOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitPhi(instruction, operation)
}

type SelectOperation struct {
	ReturnedType Type
	Value        Variable
	Consequence  Variable
	Alternative  Variable
	Output       *Field
	Block        *Block
}

type ReturnOperation struct {
	ReturnedType Type
	Value        Variable
	Block        *Block
}

func (operation *ReturnOperation) IsReturningValue() bool {
	return operation.ReturnedType.Matches(VoidType)
}

func (operation *ReturnOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitReturn(instruction, operation)
}

type CreateOperation struct {
	CreatedType Type
	Output      *Field
	Arguments   []Variable
}

func (operation *CreateOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitCreate(instruction, operation)
}

type CallOperation struct {
	Output     *Field
	ReturnType Type
	Arguments  []Variable
	Method     symbol.Reference
}

func (operation *CallOperation) HasOutput() bool {
	return operation.Output != nil
}

func (operation *CallOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitCall(instruction, operation)
}
