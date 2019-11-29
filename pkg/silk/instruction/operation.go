package instruction

type Operation int

const (
	Invalid Operation = iota
	CallOperation

	arithmeticOperationBegin
	AddOperation
	SubtractOperation
	MultiplyOperation
	DivideOperation
	arithmeticOperationEnd
	logicalBegin
	CompareOperation
	logicalEnd
	FieldOperation
	controlFlowBegin
	JumpOperation
	PhiOperation
	ReturnOperation
	controlFlowEnd
	memoryBegin
	StoreOperation
	LoadOperation
	memoryEnd
)
