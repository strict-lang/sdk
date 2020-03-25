package silk

import "gitlab.com/strict-lang/sdk/pkg/compiler/typing"

type MethodContext struct {
	OperandStack *VirtualOperandStack
	Variables    []*VirtualVariable
	Arguments    []*VirtualVariable
}

func (context *MethodContext) IndexVariable(index int) *VirtualVariable {
	return context.Variables[index]
}

type VirtualVariable struct {
	Type typing.Type
}

type VirtualOperandStack struct {
	CurrentIndex int
}

func (stack *VirtualOperandStack) Pop() VirtualVariable {
	stack.CurrentIndex--
	return VirtualVariable{}
}

func (stack *VirtualOperandStack) Push(variable VirtualVariable) {
	stack.CurrentIndex--
}
