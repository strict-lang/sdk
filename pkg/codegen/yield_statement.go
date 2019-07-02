package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

const (
	yieldListName = "_GENERATED_yield_list_"
)

func (generator *CodeGenerator) GenerateYieldStatement(statement *ast.YieldStatement) {
	generator.Emitf("%s.append(", yieldListName)
	statement.Accept(generator.generators)
	generator.Emitf("%s)")
	generator.addEpilogueGenerator(generator.returnYieldList)
}

func (generator *CodeGenerator) returnYieldList() {
	generator.Emitf("return %s;", yieldListName)
}
