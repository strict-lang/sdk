package analysis

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
	"os"
	"path/filepath"
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

func (pass *ImportPass) shouldImportFile(name string) bool {
	return !strings.HasSuffix(name, pass.currentFile+strictFileExtension)
}

func filterFiles(files []string, filter func(string) bool) (filtered []string) {
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

func (pass *ImportPass) importFiles(files []string, scope scope.MutableScope) {
	importing := NewSourceImporting(files)
	if err := importing.Import(scope); err != nil {
		pass.reportFailedImport(err)
	}
}

func listFilesInDirectory(directory string) (files []string, err error) {
	err = filepath.Walk(directory, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr == nil && isStrictFile(info) {
			files = append(files, path)
		}
		return walkErr
	})
	return files, err
}

func listFilesInDirectoryFiltered(
	directory string, filter func(string) bool) ([]string, error) {

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
