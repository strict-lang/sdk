package buildtool

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"github.com/strict-lang/sdk/pkg/compiler/backend"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	"github.com/strict-lang/sdk/pkg/compiler/report"
	"path/filepath"
	"time"
)

type Build struct {
	RootPath      string
	Configuration Configuration
	Backend backend.Backend
}

type result struct {
	error error
	diagnostics *diagnostic.Diagnostics
	lineMaps *linemap.Table
}

const sourceDirectoryName = "src"

func (build *Build) Run() (report.Report, *linemap.Table, error) {
	beginTime := time.Now().UnixNano()
	result := build.run()
	return report.Report{
		Success: result.error == nil && !containsError(result.diagnostics),
		Time: report.Time{
			Begin:      beginTime,
			Completion: time.Now().UnixNano(),
		},
		Diagnostics: TranslateDiagnostics(result.diagnostics),
	}, result.lineMaps, result.error
}

func containsError(diagnostics *diagnostic.Diagnostics) bool {
	for _, entry := range diagnostics.ListEntries() {
		if entry.Kind == &diagnostic.Error {
			return true
		}
	}
	return false
}

func (build *Build) run() result {
	config, err := build.createPackageCompilationConfig()
	if err != nil {
		return result{
			error:       err,
			diagnostics: diagnostic.Empty(),
		}
	}
	packageResult := compilePackage(config)
	return result{
		diagnostics: packageResult.diagnostics,
		lineMaps: packageResult.lineMaps,
	}
}

func (build *Build) createPackageCompilationConfig() (packageCompilationConfig, error) {
	namespaces, err := build.scanNamespaces()
	if err != nil {
		return packageCompilationConfig{}, err
	}
	return packageCompilationConfig{
		backend:    build.Backend,
		namespaces: namespaces,
		outputPath: build.selectOutputDirectory(),
	}, nil
}

func (build *Build) selectOutputDirectory() string {
	hiddenDirectory := filepath.Join(build.RootPath, ".strict")
	return filepath.Join(hiddenDirectory, "build")
}

func (build *Build) scanNamespaces() (*namespace.Table, error) {
	sourceDirectoryPath := filepath.Join(build.RootPath, sourceDirectoryName)
	rootPackageName := build.Configuration.PackageName
	fileTree, err := namespace.NewFileTree(sourceDirectoryPath, rootPackageName)
	if err != nil {
		return nil, err
	}
	return fileTree.CreateTable(), nil
}
