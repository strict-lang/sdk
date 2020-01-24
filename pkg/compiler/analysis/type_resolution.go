package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

const TypeResolutionPassId = "TypeResolutionPassId"

func init() {
	passes.Register(newTypeResolution())
}

type TypeResolution struct {
	visitor tree.Visitor
	context *passes.Context
}

func newTypeResolution() *TypeResolution {
	resolution := &TypeResolution{}
	resolution.visitor = tree.NewEmptyVisitor()
	return resolution
}

func (pass *TypeResolution) Run(context *passes.Context) {
	pass.context = context
	context.Unit.Accept(pass.visitor)
}

func (pass *TypeResolution) Id() passes.Id {
	return TypeResolutionPassId
}

func (pass *TypeResolution) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, NameResolutionPassId)
}

func (pass *TypeResolution) visitStringLiteral(string *tree.StringLiteral) {
	if !isResolved(string) {
		string.ResolveType(scope.Builtins.String)
	}
}

func (pass *TypeResolution) visitNumberLiteral(number *tree.NumberLiteral) {
	if !isResolved(number) {
		pass.resolveNumberLiteral(number)
	}
}

func (pass *TypeResolution) resolveNumberLiteral(number *tree.NumberLiteral) {
	if number.IsFloat() {
		number.ResolveType(scope.Builtins.Float)
	} else {
		number.ResolveType(scope.Builtins.Number)
	}
}

func (pass *TypeResolution) resolveExpression(expression tree.Expression) *scope.Class {
	expression.Accept(pass.visitor)
	if class, ok := expression.ResolvedType(); ok {
		return class
	}
	pass.reportFailedInference(expression)
	return nil
}

func (pass *TypeResolution) resolveBinaryExpression(binary *tree.BinaryExpression) {
	if isResolved(binary) {
		return
	}
	if operation, ok := binaryOperationTypes[binary.Operator]; ok {
		leftOperandType := pass.resolveExpression(binary.LeftOperand)
		binary.ResolveType(operation(leftOperandType))
	}
	pass.reportFailedInference(binary)
}

func (pass *TypeResolution) resolveLetExpression(binding *tree.LetBinding) {
	expressionClass := pass.resolveExpression(binding.Expression)
	binding.ResolveType(expressionClass)
}

func (pass *TypeResolution) resolveUnaryExpression(unary *tree.UnaryExpression) {
	if isResolved(unary) {
		return
	}
	if operation, ok := unaryOperationTypes[unary.Operator]; ok {
		operandType := pass.resolveExpression(unary.Operand)
		unary.ResolveType(operation(operandType))
	}
	pass.reportFailedInference(unary)
}

func (pass *TypeResolution) resolveCallExpression(call *tree.CallExpression) {
	if isResolved(call) {
		return
	}
}

func (pass *TypeResolution) reportFailedInference(node tree.Node) {
	panic("Failed to infer node")
}

func isResolved(expression tree.Expression) bool {
	_, isResolved := expression.ResolvedType()
	return isResolved
}

type typeOperation func(*scope.Class) *scope.Class

func identityTypeOperation(input *scope.Class) *scope.Class {
	return input
}

func fixedTypeOperation(constantOutput *scope.Class) typeOperation {
	return func(class *scope.Class) *scope.Class {
		return constantOutput
	}
}

var alwaysBoolean = func(*scope.Class) *scope.Class {
	return scope.Builtins.Boolean
}

var binaryOperationTypes = map[token.Operator] typeOperation {
	token.SmallerOperator:       alwaysBoolean,
	token.GreaterOperator:       alwaysBoolean,
	token.EqualsOperator:        alwaysBoolean,
	token.NotEqualsOperator:     alwaysBoolean,
	token.SmallerEqualsOperator: alwaysBoolean,
	token.GreaterEqualsOperator: alwaysBoolean,
	token.AddOperator:           identityTypeOperation,
	token.SubOperator:           identityTypeOperation,
	token.MulOperator:           identityTypeOperation,
	token.DivOperator:           identityTypeOperation,
	token.ModOperator:           identityTypeOperation,
}

var unaryOperationTypes = map[token.Operator] typeOperation {
	token.NegateOperator: alwaysBoolean,
}