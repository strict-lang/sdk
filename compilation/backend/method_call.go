package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/token"
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

func (generation *Generation) GenerateCallExpression(call *syntaxtree.CallExpression) {
	if typeName, isConstructorCall := call.Method.(syntaxtree.TypeName); isConstructorCall {
		generation.generateConstructorCall(typeName, call)
		return
	}
	// TODO: Implement argument reordering by their name, if they have one.
	//  this can't be done until enough information about the called method
	//  exists. Currently, we just ignore the argument names for normal calls.
	renameBuiltinMethodNameForCall(call.Method)
	generation.EmitNode(call.Method)
	generation.emitArgumentList(call.Arguments)
}

func (generation *Generation) generateConstructorCall(typeName syntaxtree.TypeName, call *syntaxtree.CallExpression) {
	labeled, others := filterLabeledArguments(call.Arguments)
	if len(labeled) == 0 {
		generation.EmitNode(call.Method)
		generation.emitArgumentList(others)
	} else {
		generation.generateLabeledConstructorCall(typeName, labeled, others)
	}
}

func (generation *Generation) generateLabeledConstructorCall(
	typeName syntaxtree.TypeName, labeled []*syntaxtree.CallArgument, others []*syntaxtree.CallArgument) {

	generation.Emit("[&]() -> ")
	generation.EmitNode(typeName)
	generation.Emit(" {")
	generation.EmitEndOfLine()
	generation.IncreaseIndent()

	generation.DecreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
	generation.Emit("}()")

}

const constructedFieldName = "$_constructed"

func (generation *Generation) generateLabeledConstructorCallLambdaBody(
	typeName syntaxtree.TypeName, labeled []*syntaxtree.CallArgument, others []*syntaxtree.CallArgument) {

	generation.EmitNode(&syntaxtree.AssignStatement{
		Target:       &syntaxtree.FieldDeclaration{
			Name:         &syntaxtree.Identifier{
				Value:        constructedFieldName,
				NodePosition: syntaxtree.ZeroPosition{},
			},
			TypeName:     typeName,
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Value:        &syntaxtree.CallExpression{
			Method:       typeName,
			Arguments:    others,
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Operator:     token.AssignOperator,
		NodePosition: syntaxtree.ZeroPosition{},
	})
	for _, argument := range labeled {
		generation.emitFieldInitializationForLabeledArgument(argument)
	}
}

func (generation *Generation) emitFieldInitializationForLabeledArgument(argument *syntaxtree.CallArgument) {
	generation.EmitNode(&syntaxtree.AssignStatement{
		Target:       &syntaxtree.SelectExpression{
			Target:       &syntaxtree.Identifier{
				Value:        constructedFieldName,
				NodePosition: syntaxtree.ZeroPosition{},
			},
			Selection:    &syntaxtree.Identifier{
				Value:        argument.Label,
				NodePosition: syntaxtree.ZeroPosition{},
			},
			NodePosition: nil,
		},
		Value:        argument.Value,
		Operator:     token.AssignOperator,
		NodePosition: syntaxtree.ZeroPosition{},
	})
}

func filterLabeledArguments(arguments []*syntaxtree.CallArgument) (labeled []*syntaxtree.CallArgument, others []*syntaxtree.CallArgument) {
	for _, argument := range arguments {
		if argument.IsLabeled() {
			labeled = append(labeled, argument)
		} else {
			others = append(others, argument)
		}
	}
	return
}

func (generation *Generation) emitArgumentList(arguments []*syntaxtree.CallArgument) {
	generation.Emit("(")
	for index, argument := range arguments{
		if index != 0 {
			generation.Emit(", ")
		}
		generation.EmitNode(argument.Value)
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
