package namespace

import (
	"os"
	"path/filepath"
	"strings"
)

type kind int8

const (
	directoryKind kind = iota
	fileKind
)

type fileTree struct {
	entries map[string]kind
}

func (tree *fileTree) findByQualifier(qualifier string) (kind, bool) {
	kind, ok := tree.entries[qualifier]
	return kind, ok
}

func newFileTree(root string, rootNamespaceName string) (*fileTree, error) {
	entries, err := scanFiles(root, rootNamespaceName)
	if err != nil {
		return nil, err
	}
	return &fileTree{entries: entries}, nil
}

func scanFiles(
	root string,
	rootNamespaceName string) (values map[string]kind, err error) {

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		qualifier := createQualifiedName(path, rootNamespaceName)
		values[qualifier] = kindForInfo(info)
		return nil
	})
	return
}

func createQualifiedName(path string, rootNamespaceName string) string {
	pathQualifier := convertPathToQualifier(path)
	return rootNamespaceName + "." + pathQualifier
}

func convertPathToQualifier(path string) string {
	return strings.ReplaceAll(path, string(os.PathSeparator), ".")
}

func kindForInfo(info os.FileInfo) kind {
	if info.IsDir() {
		return directoryKind
	}
	return fileKind
}
