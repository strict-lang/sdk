package ilasm

type CodeBuilder struct {}

func (code *CodeBuilder) BeforeLastPush() *CodeBuilder {
	return nil
}

func (code *CodeBuilder) EmitPush(class *Class) {

}

func (code *CodeBuilder) EmitFieldLoad(variable *VirtualVariable) {

}

func (code *CodeBuilder) EmitFieldStore(variable *VirtualVariable) {

}

func (code *CodeBuilder) EmitMemberLoad(field MemberField) {

}

func (code *CodeBuilder) EmitMemberStore(field MemberField) {

}

func (code *CodeBuilder) PushStringConstant(value string) {}

func (code *CodeBuilder) PushNumberConstant(class *Class, value string) {

}

func (code *CodeBuilder) PushConstantInt(value int) {
}

func (code *CodeBuilder) EmitAdd(class *Class) {}
func (code *CodeBuilder) EmitSubtraction(class *Class) {}
func (code *CodeBuilder) EmitMultiplication(class *Class) {}
func (code *CodeBuilder) EmitDivision(class *Class) {}
