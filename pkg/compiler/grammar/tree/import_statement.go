package tree

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
)

type ImportStatement struct {
	Target ImportTarget
	Alias  *Identifier
	Region input.Region
	Parent Node
}

func (statement *ImportStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}

func (statement *ImportStatement) SetEnclosingNode(target Node) {
	statement.Parent = target
}

func (statement *ImportStatement) Accept(visitor Visitor) {
	visitor.VisitImportStatement(statement)
}

func (statement *ImportStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
}

func (statement *ImportStatement) Locate() input.Region {
	return statement.Region
}

func (statement *ImportStatement) Matches(node Node) bool {
	if target, ok := node.(*ImportStatement); ok {
		return statement.Target.Matches(target.Target) &&
			statement.matchesAlias(target)
	}
	return false
}

func (statement *ImportStatement) matchesAlias(target *ImportStatement) bool {
	if statement.Alias == nil {
		return target.Alias == nil
	}
	return statement.Alias.Matches(target.Alias)
}

// ImportTarget is the File or package that is imported with an ImportStatement.
type ImportTarget interface {
	Namespace() string
	Matches(ImportTarget) bool
}

// HasAlias returns true if the import has an alias clause. Alias clauses can
// not have empty alias values.
func (statement *ImportStatement) HasAlias() bool {
	return statement.Alias != nil && len(statement.Alias.Value) != 0
}

// ModuleName returns the name of the imported module.
func (statement *ImportStatement) ModuleName() string {
	if statement.HasAlias() {
		return statement.Alias.Value
	}
	return statement.Target.Namespace()
}

type IdentifierChainImport struct {
	Chain []string
}

func (target *IdentifierChainImport) Matches(entry ImportTarget) bool {
	if chainImport, isChainImport := entry.(*IdentifierChainImport); isChainImport {
		return isStringArrayEqual(target.Chain, chainImport.Chain)
	}
	return false
}

func isStringArrayEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index, element := range left {
		if element != right[index] {
			return false
		}
	}
	return true
}

func (target *IdentifierChainImport) FilePath() string {
	if len(target.Chain) == 0 {
		return ""
	}
	var path strings.Builder
	path.WriteRune('"')
	writePath(target.Chain, &path)
	path.WriteString("\"")
	return path.String()
}

func writePath(parts []string, builder *strings.Builder) {
	for index, element := range parts {
		if index != 0 {
			builder.WriteRune('/')
		}
		builder.WriteString(element)
	}
}

func (target *IdentifierChainImport) Namespace() string {
	// The module is imported into an anonymous namespace
	return ""
}

type FileImport struct {
	Path string
}

func (target *FileImport) Matches(entry ImportTarget) bool {
	if fileImport, isFileImport := entry.(*FileImport); isFileImport {
		return fileImport.Path == target.Path
	}
	return false
}

func (target *FileImport) FilePath() string {
	return fmt.Sprintf("%s", target.Path)
}

func (target *FileImport) Namespace() string {
	path := target.Path
	begin := findFileNameBegin(path)
	end := findFileNameEnd(path)
	return path[begin:end]
}

func findFileNameBegin(path string) int {
	if strings.Contains(path, "/") {
		return strings.LastIndex(path, "/") + 1
	}
	return 0
}

func findFileNameEnd(path string) int {
	if strings.HasSuffix(path, ".h") {
		return len(path) - 2
	}
	return len(path)
}
