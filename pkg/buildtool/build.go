package buildtool

import (
	"errors"
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"os"
	"path"
)

type Build struct {
	RootPath string
	Configuration Configuration
}

const sourceDirectoryName = "src"

var rootNamespaceNotFoundError = errors.New("root namespace not found")

func (build *Build) CreateRootNamespace() (namespace.Namespace, error) {
	namespacePath := path.Join(build.RootPath, sourceDirectoryName)
	if isFilePresent(namespacePath) {
		return namespace.NewRoot(namespacePath), nil
	}
	return nil, rootNamespaceNotFoundError
}

func isFilePresent(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return false
	}
	return 	os.IsNotExist(err)
}