package code

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/token"
	"log"
)

type TypeInference struct {
	visitor *ast.Visitor
	lastResult inferenceResult
	classType *Type
	currentScope *Scope
	currentCallInstanceType *Type
	lastIdentifier *ast.Identifier
}

type inferenceResult struct {
	type_ *Type
	success bool
}

func (inference *TypeInference) inferNode(node ast.Node) inferenceResult {
	node.Accept(inference.visitor)
	return inference.lastResult
}

func (inference *TypeInference) emitType(emitted *Type) {
	inference.lastResult = inferenceResult{
		type_: emitted,
		success: true,
	}
}

func (inference *TypeInference) emitError() {
	inference.lastResult = inferenceResult{
		type_: nil,
		success: true,
	}
}

var unaryOperationType = map[token.Operator] *Type {
	token.NegateOperator: builtinTypes.boolType,
}

func (inference *TypeInference) visitUnaryExpression(expression *ast.UnaryExpression) {
	if inferred := inference.inferNode(expression.Operand); !inferred.success {
		inference.emitError()
		return
	}
	operationType, supportedOperation := unaryOperationType[expression.Operator]
	if !supportedOperation {
		inference.emitType(operationType)
		return
	}
	inference.emitType(operationType)
}

func fixedTypeOperation(type_ *Type)  func(*Type) *Type {
	return func(*Type) *Type {
		return type_
	}
}

func identityTypeOperation() func(*Type) *Type {
	return func(type_ *Type) *Type {
		return type_
	}
}

var binaryOperationType = map[token.Operator] func(*Type) *Type {
	token.SmallerOperator: fixedTypeOperation(builtinTypes.boolType),
	token.GreaterOperator: fixedTypeOperation(builtinTypes.boolType),
	token.EqualsOperator: fixedTypeOperation(builtinTypes.boolType),
	token.NotEqualsOperator: fixedTypeOperation(builtinTypes.boolType),
	token.SmallerEqualsOperator: fixedTypeOperation(builtinTypes.boolType),
	token.GreaterEqualsOperator: fixedTypeOperation(builtinTypes.boolType),
	token.AddOperator: identityTypeOperation(),
	token.SubOperator: identityTypeOperation(),
	token.MulOperator: identityTypeOperation(),
	token.DivOperator: identityTypeOperation(),
	token.ModOperator: identityTypeOperation(),
}

func (inference *TypeInference) visitBinaryExpression(expression *ast.BinaryExpression) {
	leftHandSideType := inference.inferNode(expression.LeftOperand)
	rightHandSideType := inference.inferNode(expression.RightOperand)
	if !leftHandSideType.success || !rightHandSideType.success {
		inference.emitError()
		return
	}
	operationTypeFunc, supportedOperation := binaryOperationType[expression.Operator]
	if !supportedOperation {
		inference.emitError()
		return
	}
	inference.emitType(operationTypeFunc(leftHandSideType.type_))
}

func (inference *TypeInference) methodCallInstance() *Type {
	if inference.currentCallInstanceType == nil {
		return inference.classType
	} else {
		return inference.currentCallInstanceType
	}
}

func (inference *TypeInference) visitSelectorExpression(selector *ast.SelectorExpression) {
	if _, ok := selector.Selection.(*ast.Identifier); ok {

	}
}

func (inference *TypeInference)  visitMethodCall(call *ast.MethodCall) {
	inference.lastIdentifier =  nil
	defer func() {inference.lastIdentifier = nil}()
	if result := inference.inferNode(call.Method); !result.success {
		inference.emitError()
		return
	}
	callType := inference.methodCallInstance()
	methodName := inference.lastIdentifier
	if methodName == nil {
		log.Println("TypeInference - Could not resolve method name")
		inference.emitError()
		return
	}
	method, ok := callType.LookupMethod(methodName.Value)
	if !ok {
		inference.emitError()
		return
	}
	inference.emitType(method.ReturnType)
	inference.currentCallInstanceType = method.ReturnType
}

func (inference *TypeInference) visitNumberLiteral(literal *ast.NumberLiteral) {
	if literal.IsFloat() {
		inference.emitType(builtinTypes.floatType)
	} else {
		inference.emitType(builtinTypes.intType)
	}
}

func (inference *TypeInference) visitStringLiteral(literal *ast.StringLiteral) {
	inference.emitType(builtinTypes.stringType)
}