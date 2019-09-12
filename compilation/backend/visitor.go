package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

func CreateGenericCppVisitor(generation *Generation) *ast.Visitor {
	return &ast.Visitor{
		VisitListTypeName:         generation.GenerateListTypeName,
		VisitAssertStatement:      generation.GenerateAssertStatement,
		VisitConcreteTypeName:     generation.GenerateConcreteTypeName,
		VisitCreateExpression:     generation.GenerateCreateExpression,
		VisitDecrementStatement:   generation.GenerateDecrementStatement,
		VisitEmptyStatement:       generation.GenerateEmptyStatement,
		VisitGenericTypeName:      generation.GenerateGenericTypeName,
		VisitIncrementStatement:   generation.GenerateIncrementStatement,
		VisitInvalidStatement:     generation.GenerateInvalidStatement,
		VisitParameter:            generation.GenerateParameter,
		VisitTestStatement:        generation.GenerateTestStatement,
		VisitClassDeclaration:     generation.GenerateClassDeclaration,
		VisitMethodDeclaration:    generation.GenerateMethod,
		VisitIdentifier:           generation.GenerateIdentifier,
		VisitMethodCall:           generation.GenerateMethodCall,
		VisitStringLiteral:        generation.GenerateStringLiteral,
		VisitNumberLiteral:        generation.GenerateNumberLiteral,
		VisitYieldStatement:       generation.GenerateYieldStatement,
		VisitBlockStatement:       generation.GenerateBlockStatement,
		VisitReturnStatement:      generation.GenerateReturnStatement,
		VisitTranslationUnit:      generation.GenerateTranslationUnit,
		VisitFieldDeclaration:     generation.GenerateFieldDeclaration,
		VisitAssignStatement:      generation.GenerateAssignStatement,
		VisitUnaryExpression:      generation.GenerateUnaryExpression,
		VisitImportStatement:      generation.GenerateImportStatement,
		VisitBinaryExpression:     generation.GenerateBinaryExpression,
		VisitSelectorExpression:   generation.GenerateSelectorExpression,
		VisitExpressionStatement:  generation.GenerateExpressionStatement,
		VisitRangedLoopStatement:  generation.GenerateRangedLoopStatement,
		VisitConditionalStatement: generation.GenerateConditionalStatement,
		VisitForEachLoopStatement: generation.GenerateForEachLoopStatement,
	}
}
