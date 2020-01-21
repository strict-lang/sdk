package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

const TypeResolutionId = "TypeResolution"

type TypeResolution struct {
	visitor tree.Visitor
	context *passes.Context
}

func (pass *TypeResolution) Run(context *passes.Context) {
	pass.context = context
	context.Unit.Accept(pass.visitor)
}

func (pass *TypeResolution) Id() passes.Id {
	return TypeResolutionId
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