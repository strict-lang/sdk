package silk

type Visitor interface {
	VisitInstruction(*Instruction)
	VisitArithmetic(*Instruction, *ArithmeticOperation)
	VisitCall(*Instruction, *CallOperation)
	VisitSelect(*Instruction, *SelectOperation)
	VisitPhi(*Instruction, *PhiOperation)
	VisitLoad(*Instruction, *LoadOperation)
	VisitStore(*Instruction, *StoreOperation)
	VisitPush(*Instruction, *PushOperation)
	VisitPop(*Instruction, *PopOperation)
	VisitReturn(*Instruction, *ReturnOperation)
	VisitCreate(*Instruction, *CreateOperation)
}
