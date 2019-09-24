package code

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"log"
)

type TypeInference struct {
	visitor                 *syntaxtree2.Visitor
	lastResult              inferenceResult
	classType               *Type
	currentScope            *Scope
	currentCallInstanceType *Type
	lastIdentifier          *syntaxtree2.Identifier
}

type inferenceResult struct {
	type_   *Type
	success bool
}

func (inference *TypeInference) inferNode(node syntaxtree2.Node) inferenceResult {
	node.Accept(inference.visitor)
	return inference.lastResult
}

func (inference *TypeInference) emitType(emitted *Type) {
	inference.lastResult = inferenceResult{
		type_:   emitted,
		success: true,
	}
}

func (inference *TypeInference) emitError() {
	inference.lastResult = inferenceResult{
		type_:   nil,
		success: true,
	}
}

var unaryOperationType = map[token2.Operator]*Type{
	token2.NegateOperator: builtinTypes.boolType,
}

func (inference *TypeInference) visitUnaryExpression(expression *syntaxtree2.UnaryExpression) {
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

func fixedTypeOperation(type_ *Type) func(*Type) *Type {
	return func(*Type) *Type {
		return type_
	}
}

func identityTypeOperation() func(*Type) *Type {
	return func(type_ *Type) *Type {
		return type_
	}
}

var binaryOperationType = map[token2.Operator]func(*Type) *Type{
	token2.SmallerOperator:       fixedTypeOperation(builtinTypes.boolType),
	token2.GreaterOperator:       fixedTypeOperation(builtinTypes.boolType),
	token2.EqualsOperator:        fixedTypeOperation(builtinTypes.boolType),
	token2.NotEqualsOperator:     fixedTypeOperation(builtinTypes.boolType),
	token2.SmallerEqualsOperator: fixedTypeOperation(builtinTypes.boolType),
	token2.GreaterEqualsOperator: fixedTypeOperation(builtinTypes.boolType),
	token2.AddOperator:           identityTypeOperation(),
	token2.SubOperator:           identityTypeOperation(),
	token2.MulOperator:           identityTypeOperation(),
	token2.DivOperator:           identityTypeOperation(),
	token2.ModOperator:           identityTypeOperation(),
}

func (inference *TypeInference) visitBinaryExpression(expression *syntaxtree2.BinaryExpression) {
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

func (inference *TypeInference) visitSelectorExpression(selector *syntaxtree2.SelectExpression) {
	if _, ok := selector.Selection.(*syntaxtree2.Identifier); ok {

	}
}

func (inference *TypeInference) visitMethodCall(call *syntaxtree2.CallExpression) {
	inference.lastIdentifier = nil
	defer func() { inference.lastIdentifier = nil }()
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

func (inference *TypeInference) visitNumberLiteral(literal *syntaxtree2.NumberLiteral) {
	if literal.IsFloat() {
		inference.emitType(builtinTypes.floatType)
	} else {
		inference.emitType(builtinTypes.intType)
	}
}

func (inference *TypeInference) visitStringLiteral(literal *syntaxtree2.StringLiteral) {
	inference.emitType(builtinTypes.stringType)
}
