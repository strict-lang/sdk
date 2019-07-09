package codegen

import "gitlab.com/strict-lang/sdk/compiler/ast"

func (generator *CodeGenerator) GenerateConditionalStatement(statement *ast.ConditionalStatement) {
	generator.Emit("if (")
	generator.EmitNode(statement.Condition)
	generator.Emit(") ")
	generator.EmitNode(statement.Body)
	if statement.Else != nil {
		generator.Emit(" else ")
		generator.EmitNode(statement.Else)
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
}

func (generator *CodeGenerator) GenerateFromToLoopStatement(statement *ast.FromToLoopStatement) {
	generator.Emitf("for (auto %s = ", statement.Field.Value)
	generator.EmitNode(statement.From)
	generator.Emitf("; %s < ", statement.Field.Value)
	generator.EmitNode(statement.To)
	generator.Emitf("; %s++) ", statement.Field.Value)

	generator.EmitNode(statement.Body)
}

func (generator *CodeGenerator) GenerateForEachLoopStatement(statement *ast.ForeachLoopStatement) {
	generator.Emitf("for (auto %s : ", statement.Field.Value)
	generator.EmitNode(statement.Target)
	generator.Emit(") ")

	generator.EmitNode(statement.Body)
}

func (generator *CodeGenerator) GenerateReturnStatement(statement *ast.ReturnStatement) {
	if statement.Value == nil {
		generator.Emit("return;")
		return
	}
	generator.Emit("return ")
	generator.EmitNode(statement.Value)
	generator.Emit(";")
}

func (generator *CodeGenerator) GenerateAssignStatement(statement *ast.AssignStatement) {
	generator.Emit("auto ")
	generator.EmitNode(statement.Target)
	generator.Emitf(" = ")
	generator.EmitNode(statement.Value)
	generator.Emit(";")
}
