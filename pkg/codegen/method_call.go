package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

var builtinMethods = map[string] string {
	"log": "strict::Log",
	"logf": "strict::Logf",
}

func (generator *CodeGenerator) GenerateMethodCall(call *ast.MethodCall) {
	name := lookupMethodName(call.Name.Value)
	generator.Emit(name)
	generator.Emit("(")
	for index, argument := range call.Arguments {
		if index != 0 {
			generator.Emit(", ")
		}
		argument.Accept(generator.generators)
	}
	generator.Emit(")")
}

func lookupMethodName(name string) string {
	actualName, _ := builtinMethod(name)
	return actualName
}

func builtinMethod(name string) (string, bool) {
	builtin, ok := builtinMethods[name]
	if !ok {
		return name, false
	}
	return builtin, true
}