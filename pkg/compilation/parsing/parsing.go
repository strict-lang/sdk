package parsing

import (
	code2 "gitlab.com/strict-lang/sdk/pkg/compilation/code"
	diagnostic2 "gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

const notParsingMethod = ""

// Parsing represents the process of parsing a stream of tokens and turning them
// into an abstract syntax tree. This class is not reusable and can only produce
// one translation unit. It does some scope management but does not do to many
// checks that could be considered semantic.
type Parsing struct {
	tokenReader token2.Stream
	rootScope   *code2.Scope
	recorder    *diagnostic2.Bag
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
	Indent token2.Indent
	Scope  *code2.Scope
	Parent *Block
}

func (parsing *Parsing) parseImportStatementList() (imports []*syntaxtree2.ImportStatement, failed []syntaxtree2.Node) {
	for token2.HasKeywordValue(parsing.token(), token2.ImportKeyword) {
		result := parsing.parseImportStatement()
		if importStatement, isImport := result.(*syntaxtree2.ImportStatement); isImport {
			imports = append(imports, importStatement)
		} else {
			failed = append(failed, result)
		}
	}
	return
}

func (parsing *Parsing) parseClassDeclaration() *syntaxtree2.ClassDeclaration {
	begin := parsing.offset()
	nodes := parsing.parseTopLevelNodes()
	return &syntaxtree2.ClassDeclaration{
		Name:         parsing.unitName,
		Parameters:   []syntaxtree2.ClassParameter{},
		SuperTypes:   []syntaxtree2.TypeName{},
		Children:     nodes,
		NodePosition: parsing.createPosition(begin),
	}
}

// ParseTranslationUnit invokes the parsing on the translation unit.
// This method can only be called once on the Parsing instance.
func (parsing *Parsing) ParseTranslationUnit() (*syntaxtree2.TranslationUnit, error) {
	begin := parsing.offset()
	imports, _ := parsing.parseImportStatementList()
	class := parsing.parseClassDeclaration()
	return &syntaxtree2.TranslationUnit{
		Name:         parsing.unitName,
		Imports:      imports,
		Class:        class,
		NodePosition: parsing.createPosition(begin),
	}, nil
}

// openBlock opens a new block of code, updates the parsing block pointer and
// creates a new scope for that block that is a child-scope of the parsers
// last block. Only statements with the blocks indent may go into the block.
func (parsing *Parsing) openBlock(indent token2.Indent) {
	var blockScope *code2.Scope
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

func (parsing *Parsing) token() token2.Token {
	return parsing.tokenReader.Last()
}

func (parsing *Parsing) advance() {
	parsing.tokenReader.Pull()
	parsing.isAtBeginOfStatement = false
}

func (parsing *Parsing) peek() token2.Token {
	return parsing.tokenReader.Peek()
}

func (parsing *Parsing) closeBlock() {
	parsing.block = parsing.block.Parent
}

func (parsing *Parsing) offset() source2.Offset {
	return parsing.token().Position().Begin()
}

func (parsing *Parsing) isParsingMethod() bool {
	return parsing.currentMethodName != notParsingMethod
}

func (parsing *Parsing) parseTopLevelNodes() []syntaxtree2.Node {
	beginOffset := parsing.offset()
	block, err := parsing.parseStatementBlock()
	if err != nil {
		return []syntaxtree2.Node{parsing.createInvalidStatement(beginOffset, err)}
	}
	return block.Children
}
