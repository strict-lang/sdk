package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
)

var builtinMethods = map[string]string{
	"log":         "puts",
	"logf":        "printf",
	"inputNumber": "Strict::InputNumber",
	"input":       "Strict::Input",
	"asString":    "c_str",
}

type identifierVisitor func(identifier *syntaxtree.Identifier)

func visitMethodName(node syntaxtree.Node, visitor identifierVisitor) bool {
	if identifier, isIdentifier := node.(*syntaxtree.Identifier); isIdentifier {
		visitor(identifier)
		return true
	}
	if selection, isSelection := node.(*syntaxtree.SelectExpression); isSelection {
		last, ok := findLastSelection(selection)
		if !ok {
			return false
		}
		return visitMethodName(last, visitor)
	}
	return false
}

func findLastSelection(expression *syntaxtree.SelectExpression) (node syntaxtree.Node, ok bool) {
	if next, ok := expression.Selection.(*syntaxtree.SelectExpression); ok {
		return findLastSelection(next)
	}
	return expression.Selection, true
}

func renameBuiltinMethodName(identifier *syntaxtree.Identifier) {
	identifier.Value = lookupMethodName(identifier.Value)
}

func renameBuiltinMethodNameForCall(node syntaxtree.Node) {
	visitMethodName(node, renameBuiltinMethodName)
}

func (generation *Generation) GenerateMethodCall(call *syntaxtree.CallExpression) {
	renameBuiltinMethodNameForCall(call.Method)
	generation.EmitNode(call.Method)
	generation.Emit("(")
	for index, argument := range call.Arguments {
		if index != 0 {
			generation.Emit(", ")
		}
		generation.EmitNode(argument)
	}
	generation.Emit(")")
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
