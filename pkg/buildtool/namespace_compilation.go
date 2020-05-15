package buildtool

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"github.com/strict-lang/sdk/pkg/compiler/analysis"
	"github.com/strict-lang/sdk/pkg/compiler/analysis/entering"
	"github.com/strict-lang/sdk/pkg/compiler/analysis/semantic"
	"github.com/strict-lang/sdk/pkg/compiler/backend"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	isolates "github.com/strict-lang/sdk/pkg/compiler/isolate"
	"github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func compileNamespace(
	lineMaps *linemap.Table,
	backend backend.Backend,
	namespace namespace.Namespace,
	namespaces *namespace.Table) *diagnostic.Diagnostics {


	compilation := newNamespaceCompilation(lineMaps, backend, namespace, namespaces)
	compilation.run()
	return compilation.diagnostics
}

type namespaceCompilation struct {
	diagnostics *diagnostic.Diagnostics
	units       []*tree.TranslationUnit
	scope       scope.Scope
	namespace   namespace.Namespace
	namespaces  *namespace.Table
	backend     backend.Backend
	lineMaps    *linemap.Table
}

func newNamespaceCompilation(
	lineMaps *linemap.Table,
	backend backend.Backend,
	namespace namespace.Namespace,
	namespaces *namespace.Table) *namespaceCompilation {

	return &namespaceCompilation{
		namespace: namespace,
		namespaces: namespaces,
		diagnostics: diagnostic.Empty(),
		backend: backend,
		lineMaps: lineMaps,
	}
}

func (compilation *namespaceCompilation) run() {
	symbol := compilation.createNamespace()
	scope.GlobalNamespaceTable().Insert(symbol.QualifiedName, symbol)
	compilation.generateOutputForAll()
}

func (compilation *namespaceCompilation) generateOutputForAll() {
	for _, unit := range compilation.units {
		go compilation.generateOutputLogged(unit)
	}
}

func (compilation *namespaceCompilation) generateOutputLogged(
	unit *tree.TranslationUnit) {

	err := compilation.generateOutput(unit)
	if err != nil {
		log.Printf("failed to compile %s: %s", unit.Name, err)
	}
}

func (compilation *namespaceCompilation) generateOutput(
	unit *tree.TranslationUnit) error {

	// TODO: Report diagnostics back to shared instance.
	//  This has to be done using some kind of synchronization.
	output, err := compilation.backend.Generate(backend.Input{
		Unit:        unit,
		Diagnostics: diagnostic.NewBag(),
	})
	if err != nil {
		return err
	}
	for _, file := range output.GeneratedFiles {
		if err := file.Save(); err != nil {
			return err
		}
	}
	return nil
}

func (compilation *namespaceCompilation) createNamespace() *scope.Namespace {
	symbol := compilation.createEmptyNamespace()
	compilation.runEarlyEnteringForAll()
	compilation.completeAnalysisForAll()
	return symbol
}

func (compilation *namespaceCompilation) completeAnalysisForAll() {
	for _, unit := range compilation.units {
		compilation.completeAnalysis(unit)
	}
}

func (compilation *namespaceCompilation) completeAnalysis(unit *tree.TranslationUnit) {
	recorder := diagnostic.NewBag()
	context :=&pass.Context{
		Unit:       unit,
		Diagnostic: recorder,
		Isolate:    isolates.New(),
	}

	if err := semantic.Run(context); err != nil {
		log.Printf("could not run analysis entering: %s", err)
	}
	diagnostics := recorder.CreateDiagnostics(diagnostic.ConvertWithLineMap(unit.LineMap))
	compilation.addDiagnostics(diagnostics)
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
	if err := entering.Run(context); err != nil {
		log.Printf("could not run early entering: %s", err)
	}
	diagnostics := recorder.CreateDiagnostics(diagnostic.ConvertWithLineMap(unit.LineMap))
	compilation.addDiagnostics(diagnostics)
}

func (compilation *namespaceCompilation) prepareIsolate(
	unit *tree.TranslationUnit) *isolates.Isolate {

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
	compilation.lineMaps.Insert(result.TranslationUnit.Name, result.LineMap)
	return result.TranslationUnit, result.Error
}


func (compilation *namespaceCompilation) addDiagnostics(
	diagnostics *diagnostic.Diagnostics) {

	compilation.diagnostics = compilation.diagnostics.Merge(diagnostics)
}
