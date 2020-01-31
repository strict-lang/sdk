package analysis

import (
	"fmt"
	"os"
	"strict.dev/sdk/pkg/compiler/grammar/syntax"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/scope"
)

type SourceImporting struct {
	files []string
	scope scope.MutableScope
}

func NewSourceImporting(files []string) Importing {
	return &SourceImporting{files: files}
}

type parsedDirectory struct {
	units []*tree.TranslationUnit
}

func (importing *SourceImporting) Import(scope scope.MutableScope) error {
	importing.scope = scope
	parsedDirectory := importing.parseFiles()
	importing.runFirstPass(parsedDirectory)
	importing.runSecondPass(parsedDirectory)
	return nil
}

func (importing *SourceImporting) parseFiles() (parsed parsedDirectory) {
	results := importing.parseFilesAsync()
	for unit := range results {
		parsed.units = append(parsed.units, unit)
	}
	return parsed
}

func (importing *SourceImporting) parseFilesAsync() <- chan *tree.TranslationUnit {
	results := make(chan *tree.TranslationUnit, len(importing.files))
	for _, file := range importing.files {
		go func(file string) {
			results <- importing.parseFile(file)
		}(file)
	}
	return results
}

func (importing *SourceImporting) parseFile(name string) *tree.TranslationUnit {
	file, err := os.Open(name)
	if err != nil {
		defer file.Close()
		result := syntax.Parse(name, input.NewStreamReader(file))
		if result.Error != nil {
			return result.TranslationUnit
		}
		importing.reportFailedParse(name, result.Error)
	} else {
		importing.reportFailedParse(name, err)
	}
	return nil
}

func (importing *SourceImporting) reportFailedParse(name string, err error) {
	panic(fmt.Errorf("could not parse %s: %s", name, err.Error()))
}

func (importing *SourceImporting) runFirstPass(directory parsedDirectory) {
	for _, unit := range directory.units {
		importing.importClass(unit.Class)
	}
}

func (importing *SourceImporting) runSecondPass(directory parsedDirectory) {
	visitor := importing.createSecondPassVisitor()
	for _, unit := range directory.units {
		for _, child := range unit.Class.Children {
			child.Accept(visitor)
		}
	}
}

func (importing *SourceImporting) createSecondPassVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.FieldDeclarationVisitor = importing.importField
	visitor.MethodDeclarationVisitor = importing.importMethod
	return visitor
}

func (importing *SourceImporting) importClass(class *tree.ClassDeclaration) {
	importing.scope.Insert(&scope.Class{
		DeclarationName: class.Name,
		ActualClass:     class.NewActualClass(),
	})
}

func (importing *SourceImporting) importMethod(method *tree.MethodDeclaration) {
	if enclosingClass, ok := tree.SearchEnclosingClass(method); ok {
		classSymbol := importing.resolveClassInScope(enclosingClass.Name)
		returnClassSymbol := importing.resolveClassInScope(method.Type.BaseName())
		parameters := importing.createParameterSymbols(classSymbol, method)
		classSymbol.Scope.Insert(&scope.Method{
			DeclarationName: method.Name.Value,
			ReturnType:      returnClassSymbol,
			Parameters:      parameters,
		})
	}
}

func (importing *SourceImporting) importField(field *tree.FieldDeclaration) {
	if enclosingClass, ok := tree.SearchEnclosingClass(field); ok {
		classSymbol := importing.resolveClassInScope(enclosingClass.Name)
		fieldClassSymbol := importing.resolveClassInScope(field.TypeName.BaseName())
		classSymbol.Scope.Insert(&scope.Field{
			DeclarationName: field.Name.Value,
			Class:           fieldClassSymbol,
			Kind:            scope.MemberField,
			EnclosingClass:  classSymbol,
		})
	}
}

func (importing *SourceImporting) createParameterSymbols(
	enclosingClass *scope.Class,
	method *tree.MethodDeclaration) []*scope.Field {

	symbols := make([]*scope.Field, len(method.Parameters))
	for index, parameter := range method.Parameters {
		symbols[index] = &scope.Field{
			DeclarationName: parameter.Name.Value,
			Class:           importing.resolveClassInScope(parameter.Type.BaseName()),
			Kind:            scope.ParameterField,
			EnclosingClass:  enclosingClass,
		}
	}
	return symbols
}

func (importing *SourceImporting) resolveClassInScope(name string) *scope.Class {
	point := scope.NewReferencePoint(name)
	if symbols := importing.scope.Lookup(point); !symbols.IsEmpty() {
		if class, ok := scope.AsClassSymbol(symbols.First().Symbol); ok {
			return class
		}
	}
	importing.reportMissingClass(name)
	return nil
}

func (importing *SourceImporting) reportMissingClass(name string) {
	panic(fmt.Errorf("class %s is not known. The imported code is invalid", name))
}
