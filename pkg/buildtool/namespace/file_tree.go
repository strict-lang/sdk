package namespace

import (
	"io/ioutil"
	"os"
	"strings"
)

type kind int8

const (
	directoryKind kind = iota
	fileKind
)

type FileTree struct {
	entries map[string]treeEntry
}

type treeEntry struct {
	kind kind
	name string
	children []treeEntry
}

func (tree *FileTree) LookupNamespace(qualifier string) Namespace {
	if treeEntry, ok := tree.entries[qualifier]; ok {
		entries := translateEntries(treeEntry.children)
		return &namespace{
			name:          extractNameFromQualifiedName(qualifier),
			qualifiedName: qualifier,
			entries:       entries,
		}
	}
	return nil
}

func (tree *FileTree) CreateTable() *Table {
	namespaces := map[string] Namespace{}
	for _, namespace := range tree.ListNamespaces() {
		namespaces[namespace.QualifiedName()] = namespace
	}
	return &Table{namespaces: namespaces}
}

func (tree *FileTree) ListNamespaces() (namespaces []Namespace) {
	for qualifiedName, treeEntry := range tree.entries {
		namespace := &namespace{
			name:          extractNameFromQualifiedName(qualifiedName),
			qualifiedName: qualifiedName,
			entries:       translateEntries(treeEntry.children),
		}
		namespaces = append(namespaces, namespace)
	}
	return
}

func translateEntries(entries []treeEntry) (result []Entry) {
	for _, treeEntry := range entries {
		result = append(result, &entry{
			fileName:  treeEntry.name,
			directory: treeEntry.kind == directoryKind,
		})
	}
	return
}

func extractNameFromQualifiedName(qualifiedName string) string {
	lastDot := strings.LastIndex(qualifiedName, ".")
	if lastDot == -1 {
		return qualifiedName
	}
	return qualifiedName[lastDot:]
}

func (tree *FileTree) findByQualifier(qualifier string) (treeEntry, bool) {
	entry, ok := tree.entries[qualifier]
	return entry, ok
}

func NewFileTree(root string, rootNamespaceName string) (*FileTree, error) {
	entries, err := scanFiles(root, rootNamespaceName)
	if err != nil {
		return nil, err
	}
	return &FileTree{entries: entries}, nil
}

func scanFiles(
	root string,
	rootNamespaceName string) (map[string]treeEntry, error) {

	creation := &treeCreation{
		rootNamespace: rootNamespaceName,
		mapped: map[string]treeEntry{},
	}
	if _, err := creation.createRecursive(root); err != nil {
		return nil, err
	}
	return creation.mapped, nil
}

type treeCreation struct {
  rootNamespace string
  mapped map[string] treeEntry
}

func (creation *treeCreation) createRecursive(path string) (treeEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return treeEntry{}, err
	}
	return creation.createRecursiveWithInfo(info)
}


func (creation *treeCreation) createRecursiveWithInfo(
	info os.FileInfo) (treeEntry, error) {

	if info.IsDir() {
		return creation.createRecursiveForDirectory(info)
	}
	return creation.createFileEntry(info), nil
}

func (creation *treeCreation) createFileEntry(info os.FileInfo) treeEntry {
	qualifiedName := creation.createQualifiedName(info.Name())
	fileEntry := treeEntry{
		kind:     fileKind,
		name:     qualifiedName,
	}
	creation.mapped[qualifiedName] = fileEntry
	return fileEntry
}

func (creation *treeCreation) createRecursiveForDirectory(
	info os.FileInfo) (treeEntry, error) {

	children, err := ioutil.ReadDir(info.Name())
	if err != nil {
		return treeEntry{}, err
	}
	entries, err := creation.	createDirectoryEntries(children)
	if err != nil {
		return treeEntry{}, err
	}
	return creation.createDirectoryEntry(info, entries), nil
}

func (creation *treeCreation) createDirectoryEntry(
	info os.FileInfo,
	entries []treeEntry) treeEntry {

	qualifiedName := creation.createQualifiedName(info.Name())
	directoryEntry := treeEntry{
		kind:     directoryKind,
		name:     qualifiedName,
		children: entries,
	}
	creation.mapped[qualifiedName] = directoryEntry
	return directoryEntry
}

func (creation *treeCreation) createDirectoryEntries(
	children []os.FileInfo) ([]treeEntry, error){

	var entries []treeEntry
	for _, child := range children {
		childTree, err := creation.createRecursiveWithInfo(child)
		if err != nil {
			return []treeEntry{}, err
		}
		entries = append(entries, childTree)
	}
	return entries, nil
}

func  (creation *treeCreation) createQualifiedName(path string) string {
	pathQualifier := convertPathToQualifier(path)
	return creation.rootNamespace + "." + pathQualifier
}

func convertPathToQualifier(path string) string {
	return strings.ReplaceAll(path, string(os.PathSeparator), ".")
}
