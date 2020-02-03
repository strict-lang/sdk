package ilasm

import (
	"strict.dev/sdk/pkg/compiler/grammar/token"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/scope"
	"strict.dev/sdk/pkg/compiler/typing"
)

func resolveClassOfExpression(expression tree.Expression) *Class {
	return nil
}

func translateClass(class typing.Type) *Class {
	return nil
}

func (generation *Generation) EmitIdentifier(identifier *tree.Identifier) {
	if field, ok := scope.AsFieldSymbol(identifier.Binding()); ok {
		generation.EmitField(field)
	}
}

func (generation *Generation) EmitField(field *scope.Field) {
	switch field.Kind {
	case scope.MemberField: generation.EmitMemberField(field)
	case scope.ConstantField: generation.EmitConstantField(field)
	case scope.VariableField: generation.emitVariableFieldLoad(field)
	case scope.ParameterField: generation.emitParameterFieldLoad(field)
	}
}

func (generation *Generation) emitVariableFieldLoad(field *scope.Field) {
	generation.emitLocalFieldLoad(field)
}

func (generation *Generation) emitParameterFieldLoad(field *scope.Field) {
	generation.emitLocalFieldLoad(field)
}

func (generation *Generation) emitLocalFieldLoad(field *scope.Field) {
	class := translateClass(field.Class.ActualClass)
	location := createLocationOfField(class, field)
	if err := location.EmitLoad(generation.code); err != nil {
		panic("failed to emit variable field")
	}
}

func createLocationOfField(class *Class, field *scope.Field) StorageLocation {
	return &VariableLocation{
		Variable:  &VirtualVariable{
			Name:  field.Name(),
			Class: class,
		},
		Parameter: field.Kind == scope.ParameterField,
	}
}


func (generation *Generation) EmitConstantField(field *scope.Field) {
	// TODO: Support constant values
	panic("Constant values are not supported")
}

func (generation *Generation) EmitMemberField(field *scope.Field) {
	location := generation.createMemberFieldLocation(field)
	if err := location.EmitLoad(generation.code); err != nil {
		panic("failed to load member field")
	}
}

func (generation *Generation) createMemberFieldLocation(field *scope.Field) StorageLocation {
	valueClass := translateClass(field.Class.ActualClass)
	enclosingClass := translateClass(field.EnclosingClass.ActualClass)
	return &MemberLocation{
		Field: MemberField{
			Name:           field.Name(),
			Class:          valueClass,
			EnclosingClass: enclosingClass,
		},
		// TODO: Implement this using FieldSelectExpressions. Currently, we cant figure
		//  out which StorageLocation the targeted field has, unless it is the own class.
		InstanceLocation: createOwnReferenceLocation(generation.currentClass),
	}
}

func (generation *Generation) EmitNumberLiteral(number *tree.NumberLiteral) {
	if number.IsFloat() {
		generation.code.PushNumberConstant(Float, number.Value)
	} else {
		if constant, ok := number.AsInt(); ok {
			generation.code.PushConstantInt(constant)
		} else {
			generation.code.PushNumberConstant(Int, number.Value)
		}
	}
}

func (generation *Generation) EmitExpression(expression tree.Expression) {}

func (generation *Generation) EmitStringLiteral(string *tree.StringLiteral) {
	generation.code.PushStringConstant(string.Value)
}

func (generation *Generation) EmitBinaryExpression(binary *tree.BinaryExpression) {
	generation.EmitExpression(binary.LeftOperand)
	generation.EmitExpression(binary.RightOperand)
	class := resolveClassOfExpression(binary)
	generation.EmitBinaryOperation(binaryOperation{
		operator:	binary.Operator,
		operandClass: class,
	})
}

type binaryOperation struct {
	operator token.Operator
	operandClass *Class
}

func (generation *Generation) EmitBinaryOperation(operation binaryOperation) {
	if emitter, ok := binaryOperationEmitters[operation.operator]; ok {
		emitter(generation.code, operation)
	} else {
		panic("unsupported binary operation")
	}
}

type binaryOperationEmitter func(code *BlockBuilder, operation binaryOperation)

var binaryOperationEmitters = map[token.Operator] binaryOperationEmitter {
	token.AddOperator: func(code *BlockBuilder, operation binaryOperation) {
		code.EmitAdd(operation.operandClass)
	},
	token.SubOperator: func(code *BlockBuilder, operation binaryOperation) {
		code.EmitSubtraction(operation.operandClass)
	},
	token.MulOperator: func(code *BlockBuilder, operation binaryOperation) {
		code.EmitMultiplication(operation.operandClass)
	},
	token.DivOperator: func(code *BlockBuilder, operation binaryOperation) {
		code.EmitDivision(operation.operandClass)
	},
}

