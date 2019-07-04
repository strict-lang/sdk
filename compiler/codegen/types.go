package codegen

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/ast"
)

const (
	BuiltinTypeText   = "std::string"
	BuiltinTypeNumber = "double"
	BuiltinTypeList   = "std::vector"
)

var builtinTypes = map[string]string{
	"text":   BuiltinTypeText,
	"number": BuiltinTypeNumber,
	"list":   BuiltinTypeList,
}

func updateGenericTypeName(name *ast.GenericTypeName) ast.TypeName {
	return ast.GenericTypeName{
		Name:    lookupTypeName(name.Name),
		Generic: updateTypeName(name.Generic),
	}
}

func updateConcreteTypeName(name *ast.ConcreteTypeName) ast.TypeName {
	return ast.ConcreteTypeName{
		Name: lookupTypeName(name.Name),
	}
}

func updateTypeName(name ast.TypeName) ast.TypeName {
	switch concrete := name.(type) {
	case *ast.GenericTypeName:
		return updateGenericTypeName(concrete)
	case *ast.ConcreteTypeName:
		return updateConcreteTypeName(concrete)
	}
	return name
}

func lookupTypeName(name string) string {
	builtin, ok := builtinTypes[name]
	if !ok {
		return name
	}
	fmt.Println(builtin)
	return builtin
}
