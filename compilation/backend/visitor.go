package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

func CreateGenericCppVisitor(generator *Generation) *ast.Visitor {
	return &ast.Visitor{
		VisitMethodDeclaration:    generator.GenerateMethod,
		VisitIdentifier:           generator.GenerateIdentifier,
		VisitMethodCall:           generator.GenerateMethodCall,
		VisitStringLiteral:        generator.GenerateStringLiteral,
		VisitNumberLiteral:        generator.GenerateNumberLiteral,
		VisitYieldStatement:       generator.GenerateYieldStatement,
		VisitBlockStatement:       generator.GenerateBlockStatement,
		VisitReturnStatement:      generator.GenerateReturnStatement,
		VisitTranslationUnit:      generator.GenerateTranslationUnit,
		VisitAssignStatement:      generator.GenerateAssignStatement,
		VisitUnaryExpression:      generator.GenerateUnaryExpression,
		VisitImportStatement:      generator.GenerateImportStatement,
		VisitBinaryExpression:     generator.GenerateBinaryExpression,
		VisitSelectorExpression:   generator.GenerateSelectorExpression,
		VisitExpressionStatement:  generator.GenerateExpressionStatement,
		VisitRangedLoopStatement:  generator.GenerateRangedLoopStatement,
		VisitConditionalStatement: generator.GenerateConditionalStatement,
		VisitForEachLoopStatement: generator.GenerateForEachLoopStatement,
	}
}
