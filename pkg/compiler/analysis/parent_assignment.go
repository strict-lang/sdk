package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/pass"
)

func assignParentRecursively(context *pass.Context) {
	context.Unit.AcceptRecursive(assignmentVisitor)
}

func NewParentAssignmentPass() *pass.Pass {
	return pass.NewPass().
		RunWith(assignParentRecursively).
		Create()
}

var assignmentVisitor = &tree.DelegatingVisitor{
	ParameterVisitor: func(parameter *tree.Parameter) {
		parameter.Type.SetEnclosingNode(parameter)
		parameter.Name.SetEnclosingNode(parameter)
	},
	IdentifierVisitor: func(identifier *tree.Identifier) {},
	CallArgumentVisitor: func(argument *tree.CallArgument) {
		argument.Value.SetEnclosingNode(argument)
	},
	ListTypeNameVisitor: func(name *tree.ListTypeName) {
		name.Element.SetEnclosingNode(name)
	},
	TestStatementVisitor: func(statement *tree.TestStatement) {
		statement.Child.SetEnclosingNode(statement)
	},
	StringLiteralVisitor: func(literal *tree.StringLiteral) {},
	NumberLiteralVisitor: func(literal *tree.NumberLiteral) {},
	CallExpressionVisitor: func(expression *tree.CallExpression) {
		for _, argument := range expression.Arguments {
			argument.SetEnclosingNode(expression)
		}
		expression.Target.SetEnclosingNode(expression)
	},
	EmptyStatementVisitor: func(statement *tree.EmptyStatement) {},
	WildcardNodeVisitor:   func(node *tree.WildcardNode) {},
	BreakStatementVisitor: func(statement *tree.BreakStatement) {},
	YieldStatementVisitor: func(statement *tree.YieldStatement) {
		statement.Value.SetEnclosingNode(statement)
	},
	BlockStatementVisitor: func(statement *tree.BlockStatement) {
		for _, child := range statement.Children {
			child.SetEnclosingNode(statement)
		}
	},
	AssertStatementVisitor: func(statement *tree.AssertStatement) {
		statement.Expression.SetEnclosingNode(statement)
	},
	UnaryExpressionVisitor: func(expression *tree.UnaryExpression) {
		expression.Operand.SetEnclosingNode(expression)
	},
	ImportStatementVisitor: func(statement *tree.ImportStatement) {
		if statement.Alias != nil {
			statement.Alias.SetEnclosingNode(statement)
		}
	},
	AssignStatementVisitor: func(statement *tree.AssignStatement) {
		statement.Target.SetEnclosingNode(statement)
		statement.Value.SetEnclosingNode(statement)
	},
	ReturnStatementVisitor: func(statement *tree.ReturnStatement) {
		if statement.Value != nil {
			statement.Value.SetEnclosingNode(statement)
		}
	},
	TranslationUnitVisitor: func(unit *tree.TranslationUnit) {
		unit.Class.SetEnclosingNode(unit)
	},
	CreateExpressionVisitor: func(expression *tree.CreateExpression) {
		expression.Type.SetEnclosingNode(expression)
		expression.Call.SetEnclosingNode(expression)
	},
	InvalidStatementVisitor: func(statement *tree.InvalidStatement) {},
	FieldDeclarationVisitor: func(declaration *tree.FieldDeclaration) {
		declaration.Name.SetEnclosingNode(declaration)
		declaration.TypeName.SetEnclosingNode(declaration)
	},
	PostfixExpressionVisitor: func(expression *tree.PostfixExpression) {
		expression.Operand.SetEnclosingNode(expression)
	},
	GenericTypeNameVisitor: func(name *tree.GenericTypeName) {
		name.Generic.SetEnclosingNode(name)
	},
	ConcreteTypeNameVisitor: func(name *tree.ConcreteTypeName) {},
	ClassDeclarationVisitor: func(declaration *tree.ClassDeclaration) {
		for _, child := range declaration.Children {
			child.SetEnclosingNode(declaration)
		}
		for _, parameter := range declaration.Parameters {
			parameter.SetEnclosingNode(declaration)
		}
	},
	BinaryExpressionVisitor: func(expression *tree.BinaryExpression) {
		expression.LeftOperand.SetEnclosingNode(expression)
		expression.RightOperand.SetEnclosingNode(expression)
	},
	MethodDeclarationVisitor: func(declaration *tree.MethodDeclaration) {
		for _, parameter := range declaration.Parameters {
			parameter.SetEnclosingNode(declaration)
		}
		declaration.Type.SetEnclosingNode(declaration)
		declaration.Body.SetEnclosingNode(declaration)
	},
	RangedLoopStatementVisitor: func(statement *tree.RangedLoopStatement) {
		statement.Body.SetEnclosingNode(statement)
		statement.Begin.SetEnclosingNode(statement)
		statement.End.SetEnclosingNode(statement)
		statement.Field.SetEnclosingNode(statement)
	},
	ExpressionStatementVisitor: func(statement *tree.ExpressionStatement) {
		statement.Expression.SetEnclosingNode(statement)
	},
	ForEachLoopStatementVisitor: func(statement *tree.ForEachLoopStatement) {
		statement.Body.SetEnclosingNode(statement)
		statement.Field.SetEnclosingNode(statement)
		statement.Sequence.SetEnclosingNode(statement)
	},
	ConditionalStatementVisitor: func(statement *tree.ConditionalStatement) {
		statement.Condition.SetEnclosingNode(statement)
		statement.Consequence.SetEnclosingNode(statement)
		if statement.Consequence != nil {
			statement.Consequence.SetEnclosingNode(statement)
		}
	},
	ListSelectExpressionVisitor: func(expression *tree.ListSelectExpression) {
		expression.Target.SetEnclosingNode(expression)
		expression.Index.SetEnclosingNode(expression)
	},
	FieldSelectExpressionVisitor: func(expression *tree.FieldSelectExpression) {
		expression.Target.SetEnclosingNode(expression)
		expression.Selection.SetEnclosingNode(expression)
	},
	ConstructorDeclarationVisitor: func(declaration *tree.ConstructorDeclaration) {
		declaration.Child.SetEnclosingNode(declaration)
		for _, parameter := range declaration.Parameters {
			parameter.SetEnclosingNode(declaration)
		}
	},
}
