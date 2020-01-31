package ilasm

type Label struct {
	name string
}

type BlockBuilder struct {
	Label *Label
	instructions []string
}

func (code *BlockBuilder) CreateNextBlock() *BlockBuilder {
	return nil
}

func (code *BlockBuilder) EmitValueReturn(class *Class) {}

func (code *BlockBuilder) EmitReturn() {}

func (code *BlockBuilder) EmitBranch(label *Label) {
}

func (code *BlockBuilder) EmitPop() {}

func (code *BlockBuilder) EmitBranchIfFalse(label *Label) {
}

func (code *BlockBuilder) BeforeLastPush() *BlockBuilder {
	return nil
}

func (code *BlockBuilder) EmitPush(class *Class) {

}

func (code *BlockBuilder) EmitFieldLoad(variable *VirtualVariable) {

}

func (code *BlockBuilder) EmitFieldStore(variable *VirtualVariable) {

}

func (code *BlockBuilder) EmitMemberLoad(field MemberField) {

}

func (code *BlockBuilder) EmitMemberStore(field MemberField) {

}

const loadStringInstruction = "ldstr"

func (code *BlockBuilder) PushStringConstant(value string) {
	quoted := "\"" + value + "\""
	code.emit(loadStringInstruction, quoted)
}

func (code *BlockBuilder) PushNumberConstant(class *Class, value string) {

}

func (code *BlockBuilder) PushConstantInt(value int) {
}

func (code *BlockBuilder) EmitAdd(class *Class)            {}
func (code *BlockBuilder) EmitSubtraction(class *Class)    {}
func (code *BlockBuilder) EmitMultiplication(class *Class) {}
func (code *BlockBuilder) EmitDivision(class *Class)       {}

func (code *BlockBuilder) emit(instruction string, operands ...string) {
}
