package backend

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

var builtinMethods = map[string]string{
	"log":         "puts",
	"logf":        "printf",
	"inputNumber": "Strict::InputNumber",
	"input":       "Strict::Input",
	"asString":    "c_str",
}

type identifierVisitor func(identifier *tree.Identifier)

func visitMethodName(node tree.Node, visitor identifierVisitor) bool {
	if identifier, isIdentifier := node.(*tree.Identifier); isIdentifier {
		visitor(identifier)
		return true
	}
	if selection, isSelection := node.(*tree.FieldSelectExpression); isSelection {
		last, ok := findLastSelection(selection)
		if !ok {
			return false
		}
		return visitMethodName(last, visitor)
	}
	return false
}

func findLastSelection(expression *tree.FieldSelectExpression) (node tree.Node, ok bool) {
	if next, ok := expression.Selection.(*tree.FieldSelectExpression); ok {
		return findLastSelection(next)
	}
	return expression.Selection, true
}

func renameBuiltinMethodName(identifier *tree.Identifier) {
	identifier.Value = lookupMethodName(identifier.Value)
}

func renameBuiltinMethodNameForCall(node tree.Node) {
	visitMethodName(node, renameBuiltinMethodName)
}

func (generation *Generation) GenerateCallExpression(call *tree.CallExpression) {
	if typeName, isConstructorCall := call.Target.(tree.TypeName); isConstructorCall {
		generation.generateConstructorCall(typeName, call)
		return
	}
	// TODO: Implement argument reordering by their name, if they have one.
	//  this can't be done until enough information about the called method
	//  exists. Currently, we just ignore the argument names for normal calls.
	renameBuiltinMethodNameForCall(call.Target)
	generation.EmitNode(call.Target)
	generation.emitArgumentList(call.Arguments)
}

func (generation *Generation) generateConstructorCall(typeName tree.TypeName, call *tree.CallExpression) {
	labeled, others := filterLabeledArguments(call.Arguments)
	if len(labeled) == 0 {
		generation.EmitNode(call.Target)
		generation.emitArgumentList(others)
	} else {
		generation.generateLabeledConstructorCall(typeName, labeled, others)
	}
}

func (generation *Generation) generateLabeledConstructorCall(
	typeName tree.TypeName, labeled []*tree.CallArgument, others []*tree.CallArgument) {

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
	typeName tree.TypeName, labeled []*tree.CallArgument, others []*tree.CallArgument) {
	generation.EmitIndent()
	generation.EmitNode(&tree.AssignStatement{
		Target: &tree.FieldDeclaration{
			Name: &tree.Identifier{
				Value: constructedFieldName,
			},
			TypeName: typeName,
		},
		Value: &tree.CallExpression{
			Target:    &tree.Identifier{
				Value: typeName.BaseName(),
				Region: typeName.Locate(),
			},
			Arguments: others,
		},
		Operator: token.AssignOperator,
	})
	for _, argument := range labeled {
		generation.EmitIndent()
		generation.emitFieldInitializationForLabeledArgument(argument)
	}
}

func (generation *Generation) emitFieldInitializationForLabeledArgument(argument *tree.CallArgument) {
	generation.EmitNode(&tree.AssignStatement{
		Target: &tree.FieldSelectExpression{
			Target: &tree.Identifier{
				Value: constructedFieldName,
			},
			Selection: &tree.Identifier{
				Value: argument.Label,
			},
		},
		Value:    argument.Value,
		Operator: token.AssignOperator,
	})
}

func filterLabeledArguments(arguments []*tree.CallArgument) (labeled []*tree.CallArgument, others []*tree.CallArgument) {
	for _, argument := range arguments {
		if argument.IsLabeled() {
			labeled = append(labeled, argument)
		} else {
			others = append(others, argument)
		}
	}
	return
}

func (generation *Generation) emitArgumentList(arguments []*tree.CallArgument) {
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
