package analysis

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"github.com/strict-lang/sdk/pkg/concurrent"
	"os"
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

func (importing *SourceImporting) parseFilesAsync() <-chan *tree.TranslationUnit {
	count := len(importing.files)
	results := make(chan *tree.TranslationUnit, count)
	latch := concurrent.NewLatch(count)
	for _, file := range importing.files {
		importing.parseFileAsync(file, latch, results)
	}
	latch.Wait()
	close(results)
	return results
}

func (importing *SourceImporting) parseFileAsync(
	file string, latch concurrent.Latch, output chan<- *tree.TranslationUnit) {

	go func() {
		output <- importing.parseFile(file)
		latch.CountDown()
	}()
}

func (importing *SourceImporting) parseFile(name string) *tree.TranslationUnit {
	file, err := os.Open(name)
	if err != nil {
		importing.reportFailedParse(name, err)
		return nil
	}
	defer file.Close()
	result := syntax.Parse(name, input.NewStreamReader(file))
	if result.Error != nil {
		result.Diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
		importing.reportFailedParse(name, result.Error)
		return nil
	}
	return result.TranslationUnit
}

func (importing *SourceImporting) reportFailedParse(name string, err error) {
	panic(fmt.Errorf("could not parse %s: %s", name, err))
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
