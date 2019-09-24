package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

var builtinMethods = map[string]string{
	"log":         "puts",
	"logf":        "printf",
	"inputNumber": "Strict::InputNumber",
	"input":       "Strict::Input",
	"asString":    "c_str",
}

type identifierVisitor func(identifier *syntaxtree2.Identifier)

func visitMethodName(node syntaxtree2.Node, visitor identifierVisitor) bool {
	if identifier, isIdentifier := node.(*syntaxtree2.Identifier); isIdentifier {
		visitor(identifier)
		return true
	}
	if selection, isSelection := node.(*syntaxtree2.SelectExpression); isSelection {
		last, ok := findLastSelection(selection)
		if !ok {
			return false
		}
		return visitMethodName(last, visitor)
	}
	return false
}

func findLastSelection(expression *syntaxtree2.SelectExpression) (node syntaxtree2.Node, ok bool) {
	if next, ok := expression.Selection.(*syntaxtree2.SelectExpression); ok {
		return findLastSelection(next)
	}
	return expression.Selection, true
}

func renameBuiltinMethodName(identifier *syntaxtree2.Identifier) {
	identifier.Value = lookupMethodName(identifier.Value)
}

func renameBuiltinMethodNameForCall(node syntaxtree2.Node) {
	visitMethodName(node, renameBuiltinMethodName)
}

func (generation *Generation) GenerateCallExpression(call *syntaxtree2.CallExpression) {
	if typeName, isConstructorCall := call.Method.(syntaxtree2.TypeName); isConstructorCall {
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

func (generation *Generation) generateConstructorCall(typeName syntaxtree2.TypeName, call *syntaxtree2.CallExpression) {
	labeled, others := filterLabeledArguments(call.Arguments)
	if len(labeled) == 0 {
		generation.EmitNode(call.Method)
		generation.emitArgumentList(others)
	} else {
		generation.generateLabeledConstructorCall(typeName, labeled, others)
	}
}

func (generation *Generation) generateLabeledConstructorCall(
	typeName syntaxtree2.TypeName, labeled []*syntaxtree2.CallArgument, others []*syntaxtree2.CallArgument) {

	generation.Emit("[&]() -> ")
	generation.EmitNode(typeName)
	generation.Emit(" {")
	generation.EmitEndOfLine()
	generation.IncreaseIndent()
	generation.generateLabeledConstructorCallLambdaBody(typeName, labeled, others)
	generation.DecreaseIndent()
	generation.EmitIndent()
	generation.Emit("}()")

}

const constructedFieldName = "$_constructed"

func (generation *Generation) generateLabeledConstructorCallLambdaBody(
	typeName syntaxtree2.TypeName, labeled []*syntaxtree2.CallArgument, others []*syntaxtree2.CallArgument) {
	generation.EmitIndent()
	generation.EmitNode(&syntaxtree2.AssignStatement{
		Target: &syntaxtree2.FieldDeclaration{
			Name: &syntaxtree2.Identifier{
				Value:        constructedFieldName,
				NodePosition: syntaxtree2.ZeroPosition{},
			},
			TypeName:     typeName,
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Value: &syntaxtree2.CallExpression{
			Method:       typeName,
			Arguments:    others,
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Operator:     token2.AssignOperator,
		NodePosition: syntaxtree2.ZeroPosition{},
	})
	for _, argument := range labeled {
		generation.EmitIndent()
		generation.emitFieldInitializationForLabeledArgument(argument)
	}
}

func (generation *Generation) emitFieldInitializationForLabeledArgument(argument *syntaxtree2.CallArgument) {
	generation.EmitNode(&syntaxtree2.AssignStatement{
		Target: &syntaxtree2.SelectExpression{
			Target: &syntaxtree2.Identifier{
				Value:        constructedFieldName,
				NodePosition: syntaxtree2.ZeroPosition{},
			},
			Selection: &syntaxtree2.Identifier{
				Value:        argument.Label,
				NodePosition: syntaxtree2.ZeroPosition{},
			},
			NodePosition: nil,
		},
		Value:        argument.Value,
		Operator:     token2.AssignOperator,
		NodePosition: syntaxtree2.ZeroPosition{},
	})
}

func filterLabeledArguments(arguments []*syntaxtree2.CallArgument) (labeled []*syntaxtree2.CallArgument, others []*syntaxtree2.CallArgument) {
	for _, argument := range arguments {
		if argument.IsLabeled() {
			labeled = append(labeled, argument)
		} else {
			others = append(others, argument)
		}
	}
	return
}

func (generation *Generation) emitArgumentList(arguments []*syntaxtree2.CallArgument) {
	generation.Emit("(")
	for index, argument := range arguments {
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
