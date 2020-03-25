package silk

type BlockId int

type Block struct {
	Id               BlockId
	predecessors     []*Block
	successors       []*Block
	firstInstruction *Instruction
	lastInstruction  *Instruction
}

func NewEntryBlock(id BlockId) *Block {
	return &Block{Id: id}
}

func NewChildBlock(block *Block, id BlockId) *Block {
	created := &Block{
		Id:           id,
		predecessors: []*Block{block},
	}
	block.AddSuccessor(created)
	return created
}

func (block *Block) AddPredecessor(predecessor *Block) {
	predecessor.successors = append(predecessor.successors, block)
	block.predecessors = append(block.predecessors, predecessor)
}

func (block *Block) AddSuccessor(successor *Block) {
	successor.predecessors = append(successor.predecessors, block)
	block.successors = append(block.successors, successor)
}

func (block *Block) Accept(visitor Visitor) {
	block.firstInstruction.Iterate(func(instruction *Instruction) {
		instruction.Accept(visitor)
	})
}

func (block *Block) PrependInstruction(instruction *Instruction) {
	instruction.block = block
	if currentHead := block.firstInstruction; currentHead != nil {
		block.firstInstruction.Prepend(instruction)
	} else {
		block.updateHeadAndTail(instruction)
	}
}

func (block *Block) AppendInstruction(instruction *Instruction) {
	instruction.block = block
	if currentTail := block.lastInstruction; currentTail != nil {
		currentTail.Append(instruction)
	} else {
		block.expectHeadIsNil()
		block.updateHeadAndTail(instruction)
	}
}

func (block *Block) expectHeadIsNil() {
	if block.firstInstruction != nil {
		panic("expected first instruction to be nil")
	}
}

func (block *Block) updateHeadAndTail(instruction *Instruction) {
	block.firstInstruction = instruction
	block.lastInstruction = instruction
}

func (block *Block) ReplaceInstruction(replaced *Instruction, target *Instruction) {
	block.firstInstruction.Iterate(func(instruction *Instruction) {
		if instruction.Matches(replaced) {
			instruction.ReplaceWith(target)
		}
	})
}

func (block *Block) IsEntry() bool {
	return len(block.predecessors) == 0
}

func (block *Block) IsExit() bool {
	return len(block.successors) == 0
}

func (block *Block) IsDominating(target *Block) bool {
	return target.IsDominatedBy(block)
}

func (block *Block) IsDominatedBy(possibleDominator *Block) bool {
	return isOnlyBlockInSlice(possibleDominator, block.predecessors)
}

func isOnlyBlockInSlice(block *Block, slice []*Block) bool {
	for _, element := range slice {
		if element != block {
			return false
		}
	}
	return true
}
