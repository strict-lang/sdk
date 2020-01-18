package silk

import "gitlab.com/strict-lang/sdk/pkg/silk/symbol"

type Operation interface {
	Matches(Operation) bool
	Accept(visitor Visitor, instruction *Instruction)
}

type Field struct {
	Id int
}

type WildcardOperation struct {}

func (WildcardOperation) Matches(Operation) bool {
	return true
}

func (WildcardOperation) Accept(visitor Visitor, instruction *Instruction) {}

type Variable interface {}

type ArithmeticOperation int

type ArithmeticInstruction struct {
	Output *Field
	Operands []Variable
	Operation ArithmeticOperation
}

type LoadOperation struct {}

type StoreOperation struct {}

type PushOperation struct {
	Type Type
	Variable Variable
	Block *Block
}

type PopOperation struct {
	Type Type
	Block *Block
}

type Comparison int

type CompareInstruction struct {
	Output *Field
	Comparison Comparison
	Operands []Variable
}

type JumpOperation struct {}

type PhiEntry struct {
	LastBlock BlockId
	Value Variable
}

type PhiOperation struct {
	ReturnedType Type
	Entries []PhiEntry
	Block     *Block
	Output *Field
}

func (operation *PhiOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitPhi(instruction, operation)
}

type SelectOperation struct {
	ReturnedType Type
	Value       Variable
	Consequence Variable
	Alternative Variable
	Output *Field
	Block *Block
}

type ReturnOperation struct {
	ReturnedType Type
	Value  Variable
	Block *Block
}

func (operation *ReturnOperation) IsReturningValue() bool {
	return operation.ReturnedType != VoidType
}

func (operation *ReturnOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitReturn(instruction, operation)
}

type CreateOperation struct {
	CreatedType Type
	Output *Field
	Arguments []Variable
}

func (operation *CreateOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitCreate(instruction, operation)
}

type CallOperation struct {
	Output *Field
	ReturnType Type
	Arguments []Variable
	Method symbol.Reference
}

func (operation *CallOperation) HasOutput() bool {
	return operation.Output != nil
}

func (operation *CallOperation) Accept(visitor Visitor, instruction *Instruction) {
	visitor.VisitCall(instruction, operation)
}
