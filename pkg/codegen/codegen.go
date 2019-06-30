package codegen

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"strings"
)

// CodeGenerator generates C code from an ast.
type CodeGenerator struct {
	unit   *ast.TranslationUnit
	output *strings.Builder
}

// NewCodeGenerator constructs a CodeGenerator that generates C code from
// the nodes in the passed translation-unit.
func NewCodeGenerator(unit *ast.TranslationUnit) *CodeGenerator {
	return &CodeGenerator{
		unit:   unit,
		output: &strings.Builder{},
	}
}

func (generator *CodeGenerator) String() string {
	return generator.output.String()
}

func (generator *CodeGenerator) Emit(code string) {
	generator.output.WriteString(code)
}

func (generator *CodeGenerator) Emitf(code string, arguments ...interface{}) {
	formatted := fmt.Sprintf(code, arguments)
	generator.output.WriteString(formatted)
}
