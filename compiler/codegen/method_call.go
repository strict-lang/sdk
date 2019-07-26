package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
)

var builtinMethods = map[string]string{
	"log":         "puts",
	"logf":        "printf",
	"inputNumber": "strict::InputNumber",
	"input": "strict::Input",
	"asString": "c_str",
}

func findMethodName(node ast.Node) (name string, ok bool){
	if identifier, isIdentifier := node.(*ast.Identifier); isIdentifier {
		name, ok = identifier.Value, true
		return
	}
	if selection, isSelection := node.(*ast.SelectorExpression); isSelection {
		last, ok := findLastSelection(selection)
		if !ok {
			return "", false
		}
		return findMethodName(last)
	}
	return "", false
}

func findLastSelection(expression *ast.SelectorExpression) (node ast.Node, ok bool) {
	if next, ok := expression.Selection.(*ast.SelectorExpression); ok {
		return findLastSelection(next)
	}
	return expression.Selection, false
}

func (generator *CodeGenerator) GenerateMethodCall(call *ast.MethodCall) {
	if identifier, ok := findMethodName(call.Method); ok {
		call.Method = &ast.Identifier{
			Value: lookupMethodName(identifier),
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
