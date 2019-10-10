package backend

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func CreateGenericCppVisitor(generation *Generation) *tree.Visitor {
	return &tree.Visitor{
		VisitListTypeName:           generation.GenerateListTypeName,
		VisitAssertStatement:        generation.GenerateAssertStatement,
		VisitConcreteTypeName:       generation.GenerateConcreteTypeName,
		VisitCreateExpression:       generation.GenerateCreateExpression,
		VisitDecrementStatement:     generation.GenerateDecrementStatement,
		VisitEmptyStatement:         generation.GenerateEmptyStatement,
		VisitGenericTypeName:        generation.GenerateGenericTypeName,
		VisitIncrementStatement:     generation.GenerateIncrementStatement,
		VisitInvalidStatement:       generation.GenerateInvalidStatement,
		VisitParameter:              generation.GenerateParameter,
		VisitTestStatement:          generation.GenerateTestStatement,
		VisitClassDeclaration:       generation.GenerateClassDeclaration,
		VisitMethodDeclaration:      generation.GenerateMethod,
		VisitIdentifier:             generation.GenerateIdentifier,
		VisitCallArgument:           func(*tree.CallArgument) {},
		VisitCallExpression:         generation.GenerateCallExpression,
		VisitStringLiteral:          generation.GenerateStringLiteral,
		VisitNumberLiteral:          generation.GenerateNumberLiteral,
		VisitYieldStatement:         generation.GenerateYieldStatement,
		VisitBlockStatement:         generation.GenerateBlockStatement,
		VisitReturnStatement:        generation.GenerateReturnStatement,
		VisitTranslationUnit:        generation.GenerateTranslationUnit,
		VisitFieldDeclaration:       generation.GenerateFieldDeclaration,
		VisitAssignStatement:        generation.GenerateAssignStatement,
		VisitUnaryExpression:        generation.GenerateUnaryExpression,
		VisitImportStatement:        generation.GenerateImportStatement,
		VisitBinaryExpression:       generation.GenerateBinaryExpression,
		VisitSelectorExpression:     generation.GenerateSelectExpression,
		VisitExpressionStatement:    generation.GenerateExpressionStatement,
		VisitRangedLoopStatement:    generation.GenerateRangedLoopStatement,
		VisitConditionalStatement:   generation.GenerateConditionalStatement,
		VisitForEachLoopStatement:   generation.GenerateForEachLoopStatement,
		VisitListSelectExpression:   generation.GenerateListSelectExpression,
		VisitConstructorDeclaration: func(declaration *tree.ConstructorDeclaration) {},
	}
}
