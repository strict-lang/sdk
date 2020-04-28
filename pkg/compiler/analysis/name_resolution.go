package analysis

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
)

const NameResolutionPassId = "NameResolutionPass"

func init() {
	registerPassInstance(&NameResolutionPass{})
}

type NameResolutionPass struct {
	context *passes.Context
	visitor tree.Visitor
}

func (pass *NameResolutionPass) Run(context *passes.Context) {
	pass.context = context
	pass.visitor = pass.createVisitor()
	context.Unit.AcceptRecursive(pass.visitor)
}

func (pass *NameResolutionPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ScopeCreationPassId, SymbolEnterPassId)
}

func (pass *NameResolutionPass) Id() passes.Id {
	return NameResolutionPassId
}

func (pass *NameResolutionPass) createVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.IdentifierVisitor = pass.visitIdentifier
	visitor.CallExpressionVisitor = pass.visitCallExpression
	visitor.StringLiteralVisitor = pass.visitStringLiteral
	visitor.NumberLiteralVisitor = pass.visitNumberLiteral
	visitor.BinaryExpressionVisitor = pass.visitBinaryExpression
	visitor.UnaryExpressionVisitor = pass.visitUnaryExpression
	visitor.LetBindingVisitor = pass.visitLetExpression
	visitor.ForEachLoopStatementVisitor = pass.visitForEachLoop
	visitor.RangedLoopStatementVisitor = pass.visitRangedLoop
	return visitor
}

func (pass *NameResolutionPass) visitIdentifier(identifier *tree.Identifier) {
	if !identifier.IsBound() && !identifier.IsPartOfDeclaration() {
		pass.resolveIdentifier(identifier)
	}
}

func (pass *NameResolutionPass) resolveIdentifier(identifier *tree.Identifier) {
	searchScope := pass.selectResolutionScope(identifier)
	if entries := searchScope.Lookup(identifier.ReferencePoint()); !entries.IsEmpty() {
		symbol := entries.First().Symbol
		identifier.Bind(symbol)
		identifier.ResolveType(pass.resolveFieldSymbolType(symbol))
	} else {
		identifier.ResolveType(scope.Builtins.Any)
		pass.reportUnresolvedField(identifier)
	}
}

func (pass *NameResolutionPass) resolveFieldSymbolType(symbol scope.Symbol) *scope.Class {
	if field, isField := scope.AsFieldSymbol(symbol); isField {
		return field.Class
	}
	if class, isClass := scope.AsClassSymbol(symbol); isClass {
		return class.ToTopLevelClassType()
	}
	if _, isMethod := scope.AsMethodSymbol(symbol); isMethod {
		return scope.TopLevelMethodType()
	}
	pass.reportUnexpectedSymbol(symbol)
	return scope.Builtins.Any
}

func (pass *NameResolutionPass) reportUnexpectedSymbol(symbol scope.Symbol) {
	log.Printf("did not expect symbol %s", symbol.Name())
}

func (pass *NameResolutionPass) visitCallExpression(call *tree.CallExpression) {
	if !isResolved(call) {
		pass.resolveCallExpression(call)
	}
}

func (pass *NameResolutionPass) resolveCallExpression(call *tree.CallExpression) {
	if name, ok := call.TargetName(); ok && !name.IsBound() {
		searchScope := pass.selectResolutionScope(call)
		if entries := searchScope.Lookup(name.ReferencePoint()); !entries.IsEmpty() {
			if methodSymbol, ok := scope.AsMethodSymbol(entries.First().Symbol); ok {
				name.Bind(methodSymbol)
				name.ResolveType(methodSymbol.ReturnType)
				call.ResolveType(methodSymbol.ReturnType)
				return
			}
		}
	}
	pass.resolveUnresolvedCall(call)
}

func (pass *NameResolutionPass) resolveUnresolvedCall(call *tree.CallExpression) {
	log.Print("could not resolve call")
	pass.context.Diagnostic.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SemanticAnalysis,
		Message:  "could not resolve call",
		UnitName: pass.context.Unit.Name,
		Error:    nil,
		Position: call.Locate(),
	})
}

func (pass *NameResolutionPass) selectResolutionScope(node tree.Expression) scope.Scope {
	if chain, ok := tree.SearchEnclosingChain(node); ok {
		index := findIndexInChain(node.Locate().Begin(), chain)
		formerIndex := index - 1
		if formerIndex >= 0 && formerIndex < len(chain.Expressions) {
			if lastType, ok := chain.Expressions[formerIndex].ResolvedType(); ok {
				return lastType.Scope
			}
			return scope.NewEmptyScope("invalid")
		}
	}
	if localScope, ok := tree.ResolveNearestScope(node); ok {
		return localScope
	}
	return scope.NewEmptyScope("invalid")
}

func findIndexInChain(position input.Offset, chain *tree.ChainExpression) int {
	for index, element := range chain.Expressions {
		if element.Locate().Begin() == position {
			return index
		}
	}
	return 0
}

func (pass *NameResolutionPass) reportUnresolvedField(identifier *tree.Identifier) {
	log.Printf("could not resolve field %s", identifier.Value)
	pass.context.Diagnostic.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SemanticAnalysis,
		Message:  "could not resolve identifier " + identifier.Value,
		UnitName: pass.context.Unit.Name,
		Error:    nil,
		Position: identifier.Locate(),
	})
}

func (pass *NameResolutionPass) visitStringLiteral(string *tree.StringLiteral) {
	if !isResolved(string) {
		string.ResolveType(scope.Builtins.String)
	}
}

func (pass *NameResolutionPass) visitNumberLiteral(number *tree.NumberLiteral) {
	if !isResolved(number) {
		pass.resolveNumberLiteral(number)
	}
}

func (pass *NameResolutionPass) resolveNumberLiteral(number *tree.NumberLiteral) {
	if number.IsFloat() {
		number.ResolveType(scope.Builtins.Float)
	} else {
		number.ResolveType(scope.Builtins.Number)
	}
}

func (pass *NameResolutionPass) resolveExpression(expression tree.Expression) *scope.Class {
	expression.Accept(pass.visitor)
	if class, ok := expression.ResolvedType(); ok {
		return class
	}
	pass.reportFailedInference(expression)
	return nil
}

func (pass *NameResolutionPass) visitBinaryExpression(binary *tree.BinaryExpression) {
	if isResolved(binary) {
		return
	}
	leftType := pass.resolveExpression(binary.LeftOperand)
	pass.resolveExpression(binary.RightOperand)
	if operation, ok := binaryOperationTypes[binary.Operator]; ok {
		binary.ResolveType(operation(leftType))
		return
	}
	binary.ResolveType(scope.Builtins.Any)
	pass.reportFailedInference(binary)
}

func (pass *NameResolutionPass) visitLetExpression(binding *tree.LetBinding) {
	expressionClass := pass.resolveExpression(binding.Expression)
	binding.ResolveType(expressionClass)
}

func (pass *NameResolutionPass) visitForEachLoop(loop *tree.ForEachLoopStatement) {
	sequenceClass := pass.resolveExpression(loop.Sequence)
	loop.Field.ResolveType(sequenceClass)
}

func (pass *NameResolutionPass) visitRangedLoop(loop *tree.RangedLoopStatement) {
	indexClass := pass.resolveExpression(loop.Begin)
	loop.Field.ResolveType(indexClass)
}

func (pass *NameResolutionPass) visitUnaryExpression(unary *tree.UnaryExpression) {
	if isResolved(unary) {
		return
	}
	pass.resolveExpression(unary.Operand)
	if operation, ok := unaryOperationTypes[unary.Operator]; ok {
		operandType := pass.resolveExpression(unary.Operand)
		unary.ResolveType(operation(operandType))
		return
	}
	unary.ResolveType(scope.Builtins.Any)
	pass.reportFailedInference(unary)
}

func (pass *NameResolutionPass) reportFailedInference(node tree.Node) {
	log.Print("could not infer type")
	pass.context.Diagnostic.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SemanticAnalysis,
		Message:  "failed to resolve type",
		UnitName: pass.context.Unit.Name,
		Position: node.Locate(),
	})
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

var binaryOperationTypes = map[token.Operator]typeOperation{
	token.SmallerOperator:       alwaysBoolean,
	token.GreaterOperator:       alwaysBoolean,
	token.EqualsOperator:        alwaysBoolean,
	token.NotEqualsOperator:     alwaysBoolean,
	token.SmallerEqualsOperator: alwaysBoolean,
	token.GreaterEqualsOperator: alwaysBoolean,
	token.AndOperator:           alwaysBoolean,
	token.OrOperator:            alwaysBoolean,
	token.XorOperator:           alwaysBoolean,
	token.AddOperator:           identityTypeOperation,
	token.SubOperator:           identityTypeOperation,
	token.MulOperator:           identityTypeOperation,
	token.DivOperator:           identityTypeOperation,
	token.ModOperator:           identityTypeOperation,
}

var unaryOperationTypes = map[token.Operator]typeOperation{
	token.NegateOperator: alwaysBoolean,
}
