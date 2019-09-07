package backend

import "gitlab.com/strict-lang/sdk/compilation/ast"

func (generation *Generation) GenerateConditionalStatement(statement *ast.ConditionalStatement) {
	generation.Emit("if (")
	generation.EmitNode(statement.Condition)
	generation.Emit(") ")
	generation.EmitNode(statement.Consequence)
	defer generation.EmitEndOfLine()
	if statement.Alternative != nil {
		generation.Emit(" else ")
		generation.EmitNode(statement.Alternative)
	}
}

const (
	yieldListName      = "$yield"
	yieldGeneratorName = "yield"
)

func (generation *Generation) GenerateYieldStatement(statement *ast.YieldStatement) {
	generation.method.addToPrologue(yieldGeneratorName, generation.declareYieldList)
	generation.method.addToEpilogue(yieldGeneratorName, generation.returnYieldList)

	generation.EmitFormatted("%s.push_back(", yieldListName)
	generation.EmitNode(statement.Value)
	generation.EmitFormatted(");")
}

func (generation *Generation) declareYieldList() {
	if generation.method == nil {
		panic("Yield statement outside of method")
	}
	generation.EmitNode(generation.method.declaration.Type)
	generation.EmitFormatted(" %s;", yieldListName)
}

func (generation *Generation) returnYieldList() {
	generation.Emit("\n")
	generation.EmitIndent()
	generation.EmitFormatted("return %s;", yieldListName)
	generation.EmitEndOfLine()
}

func (generation *Generation) GenerateRangedLoopStatement(statement *ast.RangedLoopStatement) {
	generation.EmitFormatted("for (auto %s = ", statement.ValueField.Value)
	generation.EmitNode(statement.InitialValue)
	generation.EmitFormatted("; %s < ", statement.ValueField.Value)
	generation.EmitNode(statement.EndValue)
	generation.EmitFormatted("; %s++) ", statement.ValueField.Value)

	generation.EmitNode(statement.Body)
	generation.EmitEndOfLine()
}

func (generation *Generation) GenerateForEachLoopStatement(statement *ast.ForEachLoopStatement) {
	generation.EmitFormatted("for (auto %s : ", statement.Field.Value)
	generation.EmitNode(statement.Sequence)
	generation.Emit(") ")

	generation.EmitNode(statement.Body)
	generation.EmitEndOfLine()
}

func (generation *Generation) GenerateReturnStatement(statement *ast.ReturnStatement) {
	if statement.Value == nil {
		generation.Emit("return;")
		return
	}
	generation.Emit("return ")
	generation.EmitNode(statement.Value)
	generation.Emit(";")
	generation.EmitEndOfLine()
}

func (generation *Generation) GenerateFieldDeclaration(declaration *ast.FieldDeclaration) {
	generation.EmitNode(declaration.TypeName)
	generation.EmitFormatted(" %s", declaration.Name.Value)
}

func (generation *Generation) GenerateAssignStatement(statement *ast.AssignStatement) {
	generation.EmitNode(statement.Target)
	generation.EmitFormatted(" = ")
	generation.EmitNode(statement.Value)
	generation.Emit(";")
	generation.EmitEndOfLine()
}

func (generation *Generation) GenerateBlockStatement(block *ast.BlockStatement) {
	generation.Emit("{\n")
	generation.IncreaseIndent()
	shouldAppendEndOfLineAtBegin := generation.appendNewLineAfterStatement
	generation.appendNewLineAfterStatement = false

	for index, child := range block.Children {
		if index != 0 {
			generation.Emit("\n")
		}
		generation.EmitIndent()
		generation.EmitNode(child)
	}
	generation.appendNewLineAfterStatement = shouldAppendEndOfLineAtBegin
	generation.DecreaseIndent()
	generation.Emit("\n")
	generation.EmitIndent()
	generation.Emit("}")
}

func (generation *Generation) GenerateExpressionStatement(statement *ast.ExpressionStatement) {
	generation.EmitNode(statement.Expression)
	generation.Emit(";")
	generation.EmitEndOfLine()
}
func (generation *Generation) GenerateAssertStatement(statement *ast.AssertStatement) {

	generation.Emit("if (!(")
	generation.EmitNode(statement.Expression)
	generation.Emit(")) {")
	generation.EmitEndOfLine()
	generation.IncreaseIndent()
	generation.EmitIndent()
	generation.EmitFormatted("throw \"%s\"", ComputeAssertionMessage(statement.Expression))
	generation.EmitEndOfLine()
	generation.DecreaseIndent()
	generation.Emit("}")
	generation.EmitEndOfLine()
}

func (generation *Generation) GenerateIncrementStatement(statement *ast.IncrementStatement) {
	generation.EmitNode(statement.Operand)
	generation.Emit("++")
}

func (generation *Generation) GenerateDecrementStatement(statement *ast.DecrementStatement) {
	generation.EmitNode(statement.Operand)
	generation.Emit("--")
}

func (generation *Generation) GenerateInvalidStatement(statement *ast.InvalidStatement) {
	generation.Emit("#error Invalid node at position")
}

func (generation *Generation) GenerateEmptyStatement(statement *ast.EmptyStatement) {}

func (generation *Generation) GenerateCreateExpression(create *ast.CreateExpression) {
	generation.EmitNode(create.Constructor)
}

func (generation *Generation) GenerateTestStatement(create *ast.TestStatement) {
	// Not Implemented
}