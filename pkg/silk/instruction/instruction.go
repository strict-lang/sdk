package instruction

import "fmt"

type Label struct{}

type Field struct{}

type Block struct {
	Label        Label
	Instructions []Instruction
}

type Instruction struct {
	Operation Operation
	Operands  []Operand
	Target    Field
	Offset    int
}

func (instruction *Instruction) IsArithmetic() bool {
	return instruction.isBetween(arithmeticOperationBegin, arithmeticOperationEnd)
}

func (instruction *Instruction) IsLogical() bool {
	return instruction.isBetween(logicalBegin, logicalEnd)
}

func (instruction *Instruction) IsMemory() bool {
	return instruction.isBetween(memoryBegin, memoryEnd)
}

func (instruction *Instruction) IsControlFlow() bool {
	return instruction.isBetween(controlFlowBegin, controlFlowEnd)
}

func (instruction *Instruction) isBetween(begin, end Operation) bool {
	return instruction.Operation >= begin && instruction.Operation <= end
}

func (instruction *Instruction) String() string {
	return fmt.Sprintf("%s(%v)", instruction.Operation, instruction.Operands)
}

type Visitor interface {
	VisitCall(*Instruction)
	VisitReturn(*Instruction)
	VisitAdd(*Instruction)
	VisitSubtract(*Instruction)
	VisitMultiply(*Instruction)
	VisitDivide(*Instruction)
	VisitJump(*Instruction)
	VisitLoad(*Instruction)
	VisitStore(*Instruction)
	VisitCompare(*Instruction)
	VisitField(*Instruction)
}
