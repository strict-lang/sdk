package syntax

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/code"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

var notParsingMethod = parsedMethod{name: `!none`}

// Parsing represents the process of grammar a stream of tokens and turning them
// into an abstract grammar tree. This class is not reusable and can only produce
// one translation unit. It does some scope management but does not do to many
// checks that could be considered semantic.
type Parsing struct {
	tokenReader token.Stream
	rootScope   *code.Scope
	recorder    *diagnostic.Bag
	block       *Block
	unitName    string
	// expressionDepth is the amount of parentheses encountered at the
	// current time. It is incremented every time the grammar looks at a
	// LeftParenOperator and decremented when it looks at a RightParenOperator.
	expressionDepth int
	// name of the method which is currently parsed. It is required by the optional
	// test statement within a method. The name is set to an empty string after a
	// method has been parsed.
	currentMethod        parsedMethod
	isAtBeginOfStatement bool
	structureStack *structureStack
}

// Block represents a nested sequence of statements that has a set indentation level.
// It helps the grammar to scanning code blocks and know where a block ends.
type Block struct {
	Indent token.Indent
	Scope  *code.Scope
	Parent *Block
}

func (parsing *Parsing) parseImportStatementList() (imports []*tree.ImportStatement) {
	for token.HasKeywordValue(parsing.token(), token.ImportKeyword) {
		imports = append(imports, parsing.parseImportStatement())
	}
	return imports
}

func (parsing *Parsing) parseClassDeclaration() *tree.ClassDeclaration {
	parsing.beginStructure(tree.ClassDeclarationNodeKind)
	nodes := parsing.parseTopLevelNodes()
	return &tree.ClassDeclaration{
		Name:       parsing.unitName,
		Parameters: []*tree.ClassParameter{},
		SuperTypes: []tree.TypeName{},
		Children:   nodes,
		Region:     parsing.completeStructure(tree.ClassDeclarationNodeKind),
	}
}

func (parsing *Parsing) Parse() (result *tree.TranslationUnit, err error) {
	defer func() {
		if failure := recover(); failure != nil{
			err = extractErrorFromPanic(failure)
		}
	}()
	return parsing.parseTranslationUnit(), nil
}

// parseTranslationUnit invokes the grammar on the translation unit.
// This method can only be called once on the Parsing instance.
func (parsing *Parsing) parseTranslationUnit() *tree.TranslationUnit {
	parsing.beginStructure(tree.TranslationUnitNodeKind)
	imports := parsing.parseImportStatementList()
	class := parsing.parseClassDeclaration()
	return &tree.TranslationUnit{
		Name:    parsing.unitName,
		Imports: imports,
		Class:   class,
		Region:  parsing.completeStructure(tree.TranslationUnitNodeKind),
	}
}

// openBlock opens a new block of code, updates the grammar block pointer and
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

func (parsing *Parsing) pullToken() token.Token {
	token := parsing.tokenReader.Last()
	parsing.tokenReader.Pull()
	return token
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

func (parsing *Parsing) offset() input.Offset {
	return parsing.token().Position().Begin()
}

func (parsing *Parsing) beginStructure(kind tree.NodeKind) {
	parsing.pushStructure(structureStackElement{
		beginOffset: parsing.offset(),
		nodeKind:    kind,
	})
}

func (parsing *Parsing) updateTopStructureKind(kind tree.NodeKind) {
	// Stack elements are values thus we can not change the top's field.
	// Instead the top element is exchanged.
	top, err := parsing.structureStack.pop()
	if err == nil {
		top.nodeKind = kind
		parsing.pushStructure(top)
	}
}

type InvalidCompletionError struct {
	Expected  tree.NodeKind
	Completed tree.NodeKind
}

func (err *InvalidCompletionError) Error() string {
	return fmt.Sprintf("tried to complete %s but completed %s", err.Expected, err.Completed)
}

func (parsing *Parsing) completeStructure(expectedKind tree.NodeKind) input.Region {
	structure, err := parsing.structureStack.pop()
	if err != nil {
		parsing.throwError(fmt.Errorf("could not complete %s: %s", expectedKind.Name(), err))
	}
	if expectedKind != tree.WildcardNodeKind && structure.nodeKind != expectedKind {
		parsing.throwError(&InvalidCompletionError{
			Expected:  expectedKind,
			Completed: structure.nodeKind,
		})
	}
	begin := structure.beginOffset
	return input.CreateRegion(begin, parsing.offset())
}

func (parsing *Parsing) createRegionOfCurrentStructure() input.Region {
	begin := parsing.peekStructure().beginOffset
	return input.CreateRegion(begin, parsing.offset())
}

func (parsing *Parsing) pushStructure(structure structureStackElement) {
	parsing.structureStack.push(structure)
}


func (parsing *Parsing) peekStructure() structureStackElement {
	return parsing.structureStack.peek()
}

func (parsing *Parsing) isParsingMethod() bool {
	return parsing.currentMethod != notParsingMethod
}

func (parsing *Parsing) parseTopLevelNodes() (nodes []tree.Node) {
	parsing.beginStructure(tree.BlockStatementNodeKind)
	block := parsing.parseStatementBlock()
	defer func() {
		if failure := recover(); failure != nil {
			err := extractErrorFromPanic(failure)
			invalid := parsing.completeInvalidStructure(err)
			nodes = []tree.Node{invalid}
		}
	}()
	parsing.completeStructure(tree.BlockStatementNodeKind)
	return convertStatementSliceToNodeSlice(block.Children)
}

func extractErrorFromPanic(value interface{}) error {
	if err, isError := value.(error); isError {
		return err
	}
	return fmt.Errorf("%s", value)
}

func convertStatementSliceToNodeSlice(statements []tree.Statement) (nodes []tree.Node) {
	for _, statement := range statements {
		nodes = append(nodes, statement)
	}
	return nodes
}
