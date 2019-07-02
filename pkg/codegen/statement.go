package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

func (generator *CodeGenerator) GenerateConditionalStatement(statement *ast.ConditionalStatement) {
	generator.Spaces()
	generator.Emit("if (")
	statement.Condition.Accept(generator.generators)
	generator.Emit(") ")
	statement.Body.Accept(generator.generators)
	if statement.Else != nil {
		generator.Emit(" else ")
		statement.Else.Accept(generator.generators)
	}
}

const (
	yieldListName = "$yield"
	yieldGeneratorName = "yield"
)

func (generator *CodeGenerator) GenerateYieldStatement(statement *ast.YieldStatement) {
	generator.method.addPrologueGenerator(yieldGeneratorName, generator.declareYieldList)
	generator.method.addEpilogueGenerator(yieldGeneratorName, generator.returnYieldList)

	generator.Emitf("%s.insert(", yieldListName)
	statement.Value.Accept(generator.generators)
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
