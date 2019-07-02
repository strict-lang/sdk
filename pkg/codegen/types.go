package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

var builtinTypes = map[string] string {
	"text": "strict::Text",
	"number": "strict::Number",
	"list": "strict::List",
}

func updateGenericTypeName(name *ast.GenericTypeName) {
	name.Name = lookupTypeName(name.Name)
	updateTypeName(name.Generic)
}

func updateConcreteTypeName(name *ast.ConcreteTypeName) {
	name.Name = lookupTypeName(name.Name)
}

func updateTypeName(name ast.TypeName) {
	switch concrete := name.(type) {
	case *ast.GenericTypeName:
		updateGenericTypeName(concrete)
	case *ast.ConcreteTypeName:
		updateConcreteTypeName(concrete)
	}
}

func lookupTypeName(name string) string {
	builtin, ok := builtinTypes[name]
	if !ok {
		return name
	}
	return builtin
}
