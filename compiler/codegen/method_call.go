package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
)

var builtinMethods = map[string]string{
	"log":         "puts",
	"logf":        "printf",
	"inputNumber": "strict::InputNumber",
}

func (generator *CodeGenerator) GenerateMethodCall(call *ast.MethodCall) {
	if identifier, ok := call.Method.(*ast.Identifier); ok {
		call.Method = &ast.Identifier{
			Value: lookupMethodName(identifier.Value),
		}
	}
	generator.EmitNode(call.Method)
	generator.Emit("(")
	for index, argument := range call.Arguments {
		if index != 0 {
			generator.Emit(", ")
		}
		generator.EmitNode(argument)
	}
	generator.Emit(")")
}

func lookupMethodName(name string) string {
	actualName, _ := builtinMethod(name)
	return actualName
}

func builtinMethod(name string) (string, bool) {
	builtin, ok := builtinMethods[name]
	if !ok {
		return name, false
	}
	return builtin, true
}
