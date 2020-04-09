package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func CreateGenericCppVisitor(generation *Generation) *tree.DelegatingVisitor {
	visitor := tree.NewEmptyVisitor()
	visitor.ListTypeNameVisitor = generation.GenerateListTypeName
	visitor.AssertStatementVisitor = generation.GenerateAssertStatement
	visitor.ConcreteTypeNameVisitor = generation.GenerateConcreteTypeName
	visitor.CreateExpressionVisitor = generation.GenerateCreateExpression
	visitor.PostfixExpressionVisitor = generation.GeneratePostfixExpression
	visitor.EmptyStatementVisitor = generation.GenerateEmptyStatement
	visitor.GenericTypeNameVisitor = generation.GenerateGenericTypeName
	visitor.InvalidStatementVisitor = generation.GenerateInvalidStatement
	visitor.ParameterVisitor = generation.GenerateParameter
	visitor.TestStatementVisitor = generation.GenerateTestStatement
	visitor.ClassDeclarationVisitor = generation.GenerateClassDeclaration
	visitor.MethodDeclarationVisitor = generation.GenerateMethod
	visitor.IdentifierVisitor = generation.GenerateIdentifier
	visitor.BreakStatementVisitor = generation.GenerateBreakStatement
	visitor.CallArgumentVisitor = func(*tree.CallArgument) {}
	visitor.CallExpressionVisitor = generation.GenerateCallExpression
	visitor.StringLiteralVisitor = generation.GenerateStringLiteral
	visitor.NumberLiteralVisitor = generation.GenerateNumberLiteral
	visitor.YieldStatementVisitor = generation.GenerateYieldStatement
	visitor.BlockStatementVisitor = generation.GenerateBlockStatement
	visitor.ReturnStatementVisitor = generation.GenerateReturnStatement
	visitor.TranslationUnitVisitor = generation.GenerateTranslationUnit
	visitor.FieldDeclarationVisitor = generation.GenerateFieldDeclaration
	visitor.AssignStatementVisitor = generation.GenerateAssignStatement
	visitor.UnaryExpressionVisitor = generation.GenerateUnaryExpression
	visitor.ImportStatementVisitor = generation.GenerateImportStatement
	visitor.BinaryExpressionVisitor = generation.GenerateBinaryExpression
	visitor.FieldSelectExpressionVisitor = generation.GenerateFieldSelectExpression
	visitor.ExpressionStatementVisitor = generation.GenerateExpressionStatement
	visitor.RangedLoopStatementVisitor = generation.GenerateRangedLoopStatement
	visitor.ConditionalStatementVisitor = generation.GenerateConditionalStatement
	visitor.ForEachLoopStatementVisitor = generation.GenerateForEachLoopStatement
	visitor.ListSelectExpressionVisitor = generation.GenerateListSelectExpression
	return visitor
}
