package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/isolate"
	passes "strict.dev/sdk/pkg/compiler/pass"
	"strict.dev/sdk/pkg/compiler/scope"
	"strings"
)

func init() {
	passes.Register(&ImportPass{})
}

const ImportPassId = "ImportPass"

type ImportPass struct {}

func (pass *ImportPass) Id() passes.Id {
	return ImportPassId
}

func (pass *ImportPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ScopeCreationPassId)
}

func (pass *ImportPass) Run(context *passes.Context) {
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
	pass.importDirectory(workingDirectory, scope)
}

func (pass *ImportPass) importDirectory(directory string, scope scope.MutableScope) {
	files, err := listFilesInDirectory(directory)
	if err != nil {
		pass.reportFailedImport(err)
		return
	}
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

const strictFileExtension = ".strict"

func isStrictFile(info os.FileInfo) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), strictFileExtension)
}

func (pass *ImportPass) reportFailedImport(err error) {
	panic(fmt.Errorf("failed to import: %s", err))
}

