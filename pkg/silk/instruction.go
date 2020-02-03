package silk

type Instruction struct {
	block *Block
	last *Instruction
	next *Instruction
	Operation Operation
}

func Instruct(operation Operation) *Instruction {
	return &Instruction{
		Operation: operation,
	}
}

func (instruction *Instruction) Matches(target *Instruction) bool {
	return instruction.Operation.Matches(target.Operation )
}

func (instruction *Instruction) Accept(visitor Visitor) {
	visitor.VisitInstruction(instruction)
	instruction.Operation.Accept(visitor, instruction)
}

func (instruction *Instruction) ReplaceWith(target *Instruction) {
	target.block = instruction.block
 instruction.appendToLast(target)
 instruction.prependToNext(target)
}

func (instruction *Instruction) removeSelf() {
	instruction.next = nil
	instruction.last = nil
}

func (instruction *Instruction) Append(target *Instruction) {
	target.block = instruction.block
	instruction.prependToNext(target)
	instruction.next = target
	target.last = instruction
}

func (instruction *Instruction) prependToNext(prepended *Instruction) {
	if currentNext := instruction.next; currentNext != nil {
		currentNext.last = prepended
		prepended.next = currentNext
	} else {
		instruction.block.lastInstruction = prepended
	}
}

func (instruction *Instruction) Prepend(target *Instruction) {
	target.block = instruction.block
	instruction.appendToLast(target)
	instruction.last = target
	target.next = instruction
}

func (instruction *Instruction) appendToLast(appended *Instruction) {
	if currentLast := instruction.last; currentLast != nil {
		currentLast.next = appended
		appended.last = currentLast
	} else {
		instruction.block.firstInstruction = appended
	}
}

func (instruction *Instruction) Iterate(visitor func(instruction *Instruction)) {
	current := instruction
	for current != nil {
		visitor(current)
		current = current.next
	}
}

func (instruction *Instruction) ComputeIndex() int {
	if instruction.last != nil {
		return instruction.last.ComputeIndex() + 1
	}
	return 0
}