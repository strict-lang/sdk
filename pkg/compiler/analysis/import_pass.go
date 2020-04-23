package analysis

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	passes.Register(&ImportPass{})
}

const ImportPassId = "ImportPass"

type ImportPass struct {
	currentFile string
}

func (pass *ImportPass) Id() passes.Id {
	return ImportPassId
}

func (pass *ImportPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ScopeCreationPassId)
}

func (pass *ImportPass) Run(context *passes.Context) {
	pass.currentFile = context.Unit.Name
	unitScope := ensureScopeIsMutable(context.Unit.Scope())
	for _, importStatement := range context.Unit.Imports {
		pass.processImport(importStatement, unitScope)
	}
	pass.importWorkingDirectory(unitScope)
}

func (pass *ImportPass) processImport(
	statement *tree.ImportStatement, scope scope.MutableScope) {

	pass.importDirectory(statement.Target.FilePath(), scope)
}

func (pass *ImportPass) importWorkingDirectory(scope scope.MutableScope) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		pass.reportFailedImport(err)
		return
	}
	files, err := listFilesInDirectoryFiltered(workingDirectory, pass.shouldImportFile)
	pass.importFiles(files, scope)
}

func (pass *ImportPass) shouldImportFile(info os.FileInfo) bool {
	name := info.Name()
	return strings.HasSuffix(name, pass.currentFile + strictFileExtension)
}

func filterFiles(
	files []os.FileInfo,
	filter func(info os.FileInfo) bool) (filtered []os.FileInfo) {

	for _, file := range files {
		if filter(file) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func (pass *ImportPass) importDirectory(directory string, scope scope.MutableScope) {
	files, err := listFilesInDirectory(directory)
	if err != nil {
		pass.reportFailedImport(err)
		return
	}
	pass.importFiles(files, scope)
}

func (pass *ImportPass) importFiles(files []os.FileInfo, scope scope.MutableScope) {
	names := convertInfosToNames(files)
	importing := NewSourceImporting(names)
	if err := importing.Import(scope); err != nil {
		pass.reportFailedImport(err)
	}
}

func convertInfosToNames(infos []os.FileInfo) (names []string) {
	for _, info := range infos {
		names = append(names, info.Name())
	}
	return
}

func listFilesInDirectory(directory string) (files []os.FileInfo, err error) {
	return ioutil.ReadDir(directory)
}

func listFilesInDirectoryFiltered(
	directory string,
	filter func(info os.FileInfo) bool) ([]os.FileInfo, error) {

	files, err := listFilesInDirectory(directory)
	if err != nil {
		return nil, err
	}
	return filterFiles(files, filter), nil
}

const strictFileExtension = ".strict"

func isStrictFile(info os.FileInfo) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), strictFileExtension)
}

func (pass *ImportPass) reportFailedImport(err error) {
	panic(fmt.Errorf("failed to import: %s", err))
}
