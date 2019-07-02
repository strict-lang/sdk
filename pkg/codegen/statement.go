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
	yieldListName = "_GENERATED_yield_list_"
)

func (generator *CodeGenerator) GenerateYieldStatement(statement *ast.YieldStatement) {
	generator.method.addPrologueGenerator(generator.declareYieldList)
	generator.Spaces()
	generator.Emitf("%s.append(", yieldListName)
	statement.Accept(generator.generators)
	generator.Emitf("%s)")
	generator.method.addEpilogueGenerator(generator.returnYieldList)
}

func (generator *CodeGenerator) declareYieldList() {
	generator.Emitf("List")
}

func (generator *CodeGenerator) returnYieldList() {
	generator.Spaces()
	generator.Emitf("return %s;", yieldListName)
}
