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

type identifierVisitor func(identifier *ast.Identifier)

func visitMethodName(node ast.Node, visitor identifierVisitor) bool {
	if identifier, isIdentifier := node.(*ast.Identifier); isIdentifier {
		visitor(identifier)
		return true
	}
	if selection, isSelection := node.(*ast.SelectorExpression); isSelection {
		last, ok := findLastSelection(selection)
		if !ok {
			return false
		}
		return visitMethodName(last, visitor)
	}
	return false
}

func findLastSelection(expression *ast.SelectorExpression) (node ast.Node, ok bool) {
	if next, ok := expression.Selection.(*ast.SelectorExpression); ok {
		return findLastSelection(next)
	}
	return expression.Selection, true
}

func renameBuiltinMethodName(identifier *ast.Identifier) {
	identifier.Value = lookupMethodName(identifier.Value)
}

func renameBuiltinMethodNameForCall(node ast.Node) {
	visitMethodName(node, renameBuiltinMethodName)
}

func (generator *CodeGenerator) GenerateMethodCall(call *ast.MethodCall) {
	renameBuiltinMethodNameForCall(call.Method)
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
