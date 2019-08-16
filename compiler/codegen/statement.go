package codegen

import "gitlab.com/strict-lang/sdk/compiler/ast"

func (generator *CodeGenerator) GenerateConditionalStatement(statement *ast.ConditionalStatement) {
	generator.Emit("if (")
	generator.EmitNode(statement.Condition)
	generator.Emit(") ")
	generator.EmitNode(statement.Consequence)
	defer generator.writeEndOfStatement()
	if statement.Alternative != nil {
		generator.Emit(" else ")
		generator.EmitNode(statement.Alternative)
	}
}

const (
	yieldListName      = "$yield"
	yieldGeneratorName = "yield"
)

func (generator *CodeGenerator) GenerateYieldStatement(statement *ast.YieldStatement) {
	generator.method.addPrologueGenerator(yieldGeneratorName, generator.declareYieldList)
	generator.method.addEpilogueGenerator(yieldGeneratorName, generator.returnYieldList)

	generator.Emitf("%s.push_back(", yieldListName)
	generator.EmitNode(statement.Value)
	generator.Emitf(");")
}

func (generator *CodeGenerator) declareYieldList() {
	if generator.method == nil {
		panic("Yield statement outside of method")
	}
	typeName := updateTypeName(generator.method.declaration.Type)
	generator.Spaces()
	generator.Emitf("%s %s;\n", typeName.FullName(), yieldListName)
}

func (generator *CodeGenerator) returnYieldList() {
	generator.Emit("\n")
	generator.Spaces()
	generator.Emitf("return %s;", yieldListName)
	generator.writeEndOfStatement()
}

func (generator *CodeGenerator) GenerateRangedLoopStatement(statement *ast.RangedLoopStatement) {
	generator.Emitf("for (auto %s = ", statement.ValueField.Value)
	generator.EmitNode(statement.InitialValue)
	generator.Emitf("; %s < ", statement.ValueField.Value)
	generator.EmitNode(statement.EndValue)
	generator.Emitf("; %s++) ", statement.ValueField.Value)

	generator.EmitNode(statement.Body)
	generator.writeEndOfStatement()
}

func (generator *CodeGenerator) GenerateForEachLoopStatement(statement *ast.ForEachLoopStatement) {
	generator.Emitf("for (auto %s : ", statement.Field.Value)
	generator.EmitNode(statement.Enumeration)
	generator.Emit(") ")

	generator.EmitNode(statement.Body)
	generator.writeEndOfStatement()
}

func (generator *CodeGenerator) GenerateReturnStatement(statement *ast.ReturnStatement) {
	if statement.Value == nil {
		generator.Emit("return;")
		return
	}
	generator.Emit("return ")
	generator.EmitNode(statement.Value)
	generator.Emit(";")
	generator.writeEndOfStatement()
}

func (generator *CodeGenerator) GenerateAssignStatement(statement *ast.AssignStatement) {
	generator.Emit("auto ")
	generator.EmitNode(statement.Target)
	generator.Emitf(" = ")
	generator.EmitNode(statement.Value)
	generator.Emit(";")
	generator.writeEndOfStatement()
}