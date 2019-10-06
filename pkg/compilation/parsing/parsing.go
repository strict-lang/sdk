package parsing

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/code"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

const notParsingMethod = ""

// Parsing represents the process of parsing a stream of tokens and turning them
// into an abstract syntax tree. This class is not reusable and can only produce
// one translation unit. It does some scope management but does not do to many
// checks that could be considered semantic.
type Parsing struct {
	tokenReader token.Stream
	rootScope   *code.Scope
	recorder    *diagnostic.Bag
	block       *Block
	unitName    string
	// expressionDepth is the amount of parentheses encountered at the
	// current time. It is incremented every time the parsing looks at a
	// LeftParenOperator and decremented when it looks at a RightParenOperator.
	expressionDepth int
	// name of the method which is currently parsed. It is required by the optional
	// test statement within a method. The name is set to an empty string after a
	// method has been parsed.
	currentMethodName    string
	isAtBeginOfStatement bool
}

// Block represents a nested sequence of statements that has a set indentation level.
// It helps the parsing to scanning code blocks and know where a block ends.
type Block struct {
	Indent token.Indent
	Scope  *code.Scope
	Parent *Block
}

func (parsing *Parsing) parseImportStatementList() (imports []*syntaxtree.ImportStatement, failed []syntaxtree.Node) {
	for token.HasKeywordValue(parsing.token(), token.ImportKeyword) {
		result := parsing.parseImportStatement()
		if importStatement, isImport := result.(*syntaxtree.ImportStatement); isImport {
			imports = append(imports, importStatement)
		} else {
			failed = append(failed, result)
		}
	}
	return
}

func (parsing *Parsing) parseClassDeclaration() *syntaxtree.ClassDeclaration {
	begin := parsing.offset()
	nodes := parsing.parseTopLevelNodes()
	return &syntaxtree.ClassDeclaration{
		Name:         parsing.unitName,
		Parameters:   []syntaxtree.ClassParameter{},
		SuperTypes:   []syntaxtree.TypeName{},
		Children:     nodes,
		NodePosition: parsing.createPosition(begin),
	}
}

// ParseTranslationUnit invokes the parsing on the translation unit.
// This method can only be called once on the Parsing instance.
func (parsing *Parsing) ParseTranslationUnit() (*syntaxtree.TranslationUnit, error) {
	begin := parsing.offset()
	imports, _ := parsing.parseImportStatementList()
	class := parsing.parseClassDeclaration()
	return &syntaxtree.TranslationUnit{
		Name:         parsing.unitName,
		Imports:      imports,
		Class:        class,
		NodePosition: parsing.createPosition(begin),
	}, nil
}

// openBlock opens a new block of code, updates the parsing block pointer and
// creates a new scope for that block that is a child-scope of the parsers
// last block. Only statements with the blocks indent may go into the block.
func (parsing *Parsing) openBlock(indent token.Indent) {
	var blockScope *code.Scope
	if parsing.block == nil {
		blockScope = parsing.rootScope.NewChild()
	} else {
		blockScope = parsing.block.Scope.NewChild()
	}
	block := &Block{
		Indent: indent,
		Scope:  blockScope,
		Parent: parsing.block,
	}
	parsing.block = block
}

func (parsing *Parsing) token() token.Token {
	return parsing.tokenReader.Last()
}

func (parsing *Parsing) advance() {
	parsing.tokenReader.Pull()
	parsing.isAtBeginOfStatement = false
}

func (parsing *Parsing) peek() token.Token {
	return parsing.tokenReader.Peek()
}

func (parsing *Parsing) closeBlock() {
	parsing.block = parsing.block.Parent
}

func (parsing *Parsing) offset() source.Offset {
	return parsing.token().Position().Begin()
}

func (parsing *Parsing) isParsingMethod() bool {
	return parsing.currentMethodName != notParsingMethod
}

func (parsing *Parsing) parseTopLevelNodes() []syntaxtree.Node {
	beginOffset := parsing.offset()
	block, err := parsing.parseStatementBlock()
	if err != nil {
		return []syntaxtree.Node{parsing.createInvalidStatement(beginOffset, err)}
	}
	return block.Children
}
