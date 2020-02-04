package silk

import "testing"

func newTestBlock() *Block {
	return &Block{
		Id: BlockId(0),
	}
}

func TestInstruction_Prepend(testing *testing.T) {
	block := newTestBlock()
	instruction := Instruct(NoOperation{})
	block.AppendInstruction(instruction)
	prepended := Instruct(WildcardOperation{})
	instruction.Prepend(prepended)
	if prepended.next != instruction {
		testing.Error("instruction was not prepended")
	}
	if instruction.last != prepended {
		testing.Error("last node of 2nd instruction not set")
	}
	if block.firstInstruction != prepended {
		testing.Error("block's first instruction not set")
	}
}

func TestInstruction_Append(testing *testing.T) {

}
