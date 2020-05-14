package buildtool

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"github.com/strict-lang/sdk/pkg/compiler/analysis"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	isolates "github.com/strict-lang/sdk/pkg/compiler/isolate"
	"github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type namespaceCompilation struct {
	diagnostics *diagnostic.Diagnostics
	units       []*tree.TranslationUnit
	scope       scope.Scope
	namespace   namespace.Namespace
	namespaces  *namespace.Table
}

func newNamespaceCompilation(
	namespace namespace.Namespace,
	namespaces namespace.Table) *namespaceCompilation {
	return &namespaceCompilation{
		namespace: namespace,
		namespaces: namespaces,
		diagnostics: diagnostic.Empty(),
	}
}

func (compilation *namespaceCompilation) createNamespace() *scope.Namespace {
}

func (compilation *namespaceCompilation) completeAnalysis(unit *tree.TranslationUnit) {

}

func (compilation *namespaceCompilation) runEarlyEnteringForAll() {
	for _, unit := range compilation.units {
		compilation.runEarlyEntering(unit)
	}
}

func (compilation *namespaceCompilation) runEarlyEntering(unit *tree.TranslationUnit) {
	recorder := diagnostic.NewBag()
	context := &pass.Context{
		Unit:       unit,
		Diagnostic: recorder,
		Isolate:    compilation.prepareIsolate(unit),
	}
	if err := analysis.RunEntering(context); err != nil {
		log.Printf("could not run early entering: %s", err)
	}
	diagnostics := recorder.CreateDiagnostics(diagnostic.ConvertWithLineMap(unit.LineMap))
	compilation.addDiagnostics(diagnostics)
}

func (compilation *namespaceCompilation) prepareIsolate(unit *tree.TranslationUnit) *isolates.Isolate {
	creation := &analysis.Creation{
		Unit:       unit,
		Namespaces: compilation.namespaces,
	}
	isolate := isolates.New()
	creation.Create().Store(isolate)
	return isolate
}

// TODO: Insert sub-namespaces. Could be done lazily
func (compilation *namespaceCompilation) createEmptyNamespace() *scope.Namespace {
	var classes []*scope.Class
	for _, unit := range compilation.units {
		class := compilation.createEmptyClassSymbol(unit)
		classes = append(classes, class)
	}
	return &scope.Namespace{
		DeclarationName: compilation.namespace.Name(),
		QualifiedName:   compilation.namespace.QualifiedName(),
		Scope:           scope.NewNamespaceScope(compilation.namespace, classes),
	}
}

func (compilation *namespaceCompilation) createEmptyClassSymbol(
	unit *tree.TranslationUnit) *scope.Class {

	return &scope.Class{
		DeclarationName: unit.Class.Name,
		QualifiedName:   compilation.addQualifierToName(unit.Class.Name),
	}
}

func (compilation *namespaceCompilation) addQualifierToName(name string) string {
	qualifier := compilation.namespace.QualifiedName()
	if qualifier == "" {
		return name
	}
	return name + "." + qualifier
}

func (compilation *namespaceCompilation) compileFiles() {
	for _, entry := range compilation.namespace.Entries() {
		if entry.IsDirectory() {
			continue
		}
		unit, err := compilation.compileFileAtPath(entry.FileName())
		if err != nil {
			log.Printf("failed to compile %s", entry.FileName())
			continue
		}
		compilation.units = append(compilation.units, unit)
	}
}

func (compilation *namespaceCompilation) compileFileAtPath(
	filePath string) (*tree.TranslationUnit, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return compilation.compileFile(filePath, file)
}

func (compilation *namespaceCompilation) compileFile(
	filePath string, file *os.File) (*tree.TranslationUnit, error) {

	name := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	result := syntax.Parse(name, input.NewStreamReader(file))
	compilation.addDiagnostics(result.Diagnostics)
	return result.TranslationUnit, result.Error
}


func (compilation *namespaceCompilation) addDiagnostics(
	diagnostics *diagnostic.Diagnostics) {

	compilation.diagnostics = compilation.diagnostics.Merge(diagnostics)
}
