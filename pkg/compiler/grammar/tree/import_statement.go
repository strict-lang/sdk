package tree

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
)

type ImportStatement struct {
	Target ImportTarget
	Alias  *Identifier
	Region input.Region
}

// ImportTarget is the File or package that is imported with an ImportStatement.
type ImportTarget interface {
	toModuleName() string
	FilePath() string
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
	return statement.Target.toModuleName()
}

func (statement *ImportStatement) Accept(visitor Visitor) {
	VisitImportStatement(statement)
}

func (statement *ImportStatement) AcceptRecursive(visitor Visitor) {
	VisitImportStatement(statement)
}

func (statement *ImportStatement) Locate() input.Region {
	return statement.Region
}

type IdentifierChainImport struct {
	Chain []string
}

func (target *IdentifierChainImport) FilePath() string {
	if len(target.Chain) == 0 {
		return ""
	}
	var path strings.Builder
	path.WriteRune('"')
	writePath(target.Chain, &path)
	path.WriteString(".h\"")
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

func (target *IdentifierChainImport) toModuleName() string {
	// The module is imported into an anonymous namespace
	return ""
}

type FileImport struct {
	Path string
}

func (target *FileImport) FilePath() string {
	return fmt.Sprintf("<%s>", target.Path)
}

func (target *FileImport) toModuleName() string {
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
