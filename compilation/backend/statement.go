package backend

import "gitlab.com/strict-lang/sdk/compilation/ast"

func (generation *Generation) GenerateConditionalStatement(statement *ast.ConditionalStatement) {
	generation.Emit("if (")
	generation.EmitNode(statement.Condition)
	generation.Emit(") ")
	generation.EmitNode(statement.Consequence)
	defer generation.writeEndOfStatement()
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
	generation.method.addPrologueGenerator(yieldGeneratorName, generation.declareYieldList)
	generation.method.addEpilogueGenerator(yieldGeneratorName, generation.returnYieldList)

	generation.Emitf("%s.push_back(", yieldListName)
	generation.EmitNode(statement.Value)
	generation.Emitf(");")
}

func (generation *Generation) declareYieldList() {
	if generation.method == nil {
		panic("Yield statement outside of method")
	}
	typeName := updateTypeName(generation.method.declaration.Type)
	generation.Spaces()
	generation.Emitf("%s %s;\n", typeName.FullName(), yieldListName)
}

func (generation *Generation) returnYieldList() {
	generation.Emit("\n")
	generation.Spaces()
	generation.Emitf("return %s;", yieldListName)
	generation.writeEndOfStatement()
}

func (generation *Generation) GenerateRangedLoopStatement(statement *ast.RangedLoopStatement) {
	generation.Emitf("for (auto %s = ", statement.ValueField.Value)
	generation.EmitNode(statement.InitialValue)
	generation.Emitf("; %s < ", statement.ValueField.Value)
	generation.EmitNode(statement.EndValue)
	generation.Emitf("; %s++) ", statement.ValueField.Value)

	generation.EmitNode(statement.Body)
	generation.writeEndOfStatement()
}

func (generation *Generation) GenerateForEachLoopStatement(statement *ast.ForEachLoopStatement) {
	generation.Emitf("for (auto %s : ", statement.Field.Value)
	generation.EmitNode(statement.Enumeration)
	generation.Emit(") ")

	generation.EmitNode(statement.Body)
	generation.writeEndOfStatement()
}

func (generation *Generation) GenerateReturnStatement(statement *ast.ReturnStatement) {
	if statement.Value == nil {
		generation.Emit("return;")
		return
	}
	generation.Emit("return ")
	generation.EmitNode(statement.Value)
	generation.Emit(";")
	generation.writeEndOfStatement()
}

func (generation *Generation) GenerateAssignStatement(statement *ast.AssignStatement) {
	generation.Emit("auto ")
	generation.EmitNode(statement.Target)
	generation.Emitf(" = ")
	generation.EmitNode(statement.Value)
	generation.Emit(";")
	generation.writeEndOfStatement()
}

func (generation *Generation) GenerateBlockStatement(block *ast.BlockStatement) {
	generation.Emit("{\n")
	generation.enterBlock()
	shouldAppendEndOfLineAtBegin := generation.appendNewLineAfterStatement
	generation.appendNewLineAfterStatement = false

	for index, child := range block.Children {
		if index != 0 {
			generation.Emit("\n")
		}
		generation.Spaces()
		generation.EmitNode(child)
	}
	generation.appendNewLineAfterStatement = shouldAppendEndOfLineAtBegin
	generation.leaveBlock()
	generation.Emit("\n")
	generation.Spaces()
	generation.Emit("}")
}

func (generation *Generation) GenerateExpressionStatement(statement *ast.ExpressionStatement) {
	generation.EmitNode(statement.Expression)
	generation.Emit(";")
	generation.writeEndOfStatement()
}