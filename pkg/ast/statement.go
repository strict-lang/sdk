package ast

import "fmt"

type Statement struct {

}

type MethodCall struct {
	Typed
	Method *Member
	Arguments []Expression
}

func (call MethodCall) Type() *Type {
	return call.Method.ValueType
}

func (call MethodCall) String() string {
	if call.Method == nil {
		return fmt.Sprintf("AnynomousMethodCall(%s)", call.Arguments)
	}
	return fmt.Sprintf("{%s %s(%s)}",
			call.Type(), call.Method.Name, call.Arguments)
}