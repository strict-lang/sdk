package ast

import (
	"os"
	"strings"
)

type ImportStatement struct {
	Target 				ImportTarget
	Alias        *Identifier
	NodePosition 	Position
}

func (statement *ImportStatement) HasAlias() bool {
	return statement.Alias != nil && statement.Alias.Value != ""
}

func (statement *ImportStatement) ModuleName() string {
	if statement.HasAlias() {
		return statement.Alias.Value
	}
	return statement.Target.toModuleName()
}

type ImportTarget interface {
	toModuleName() string
	FilePath() string
}

type IdentifierChainImport struct {
	Chain []string
}

func (target *IdentifierChainImport) FilePath() string {
	if len(target.Chain) == 0 {
		panic("IdentifierChainImport: Chain is empty")
		return ""
	}
	var path strings.Builder
	for index, element := range target.Chain {
		if index != 0 {
			path.WriteRune(os.PathSeparator)
		}
		path.WriteString(element)
	}
	path.WriteString(".h")
	return path.String()
}

func (target *IdentifierChainImport) toModuleName() string {
	if len(target.Chain) == 0 {
		panic("IdentifierChainImport: Chain is empty")
		return ""
	}
	return target.Chain[len(target.Chain) - 1]
}

type FileImport struct {
	Path string
}

func (target *FileImport) FilePath() string {
	return target.Path
}

func (target *FileImport) toModuleName() string {
	path := target.Path
	var begin = 0
	if strings.Contains(path, "/") {
		begin = strings.LastIndex(path, "/") + 1
	}
	var end int
	if strings.HasSuffix(path, ".h") {
		end = len(path) - 2
	} else {
		end = len(path)
	}
	return path[begin:end]
}

func (statement *ImportStatement) Accept(visitor *Visitor) {
	visitor.VisitImportStatement(statement)
}

func (statement *ImportStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitImportStatement(statement)
}

func (statement *ImportStatement) Position() Position {
	return statement.NodePosition
}

