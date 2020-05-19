package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (parsing *Parsing) parseImportStatementList() (imports []*tree.ImportStatement) {
	for token.HasKeywordValue(parsing.token(), token.ImportKeyword) {
		imports = append(imports, parsing.parseImportStatement())
	}
	return imports
}

func (parsing *Parsing) parseClassDeclaration() *tree.ClassDeclaration {
	parsing.beginStructure(tree.ClassDeclarationNodeKind)
	superTypes := parsing.parseImplementStatements()
	nodes := parsing.parseTopLevelNodes()
	return &tree.ClassDeclaration{
		Name:       convertFileNameToClassName(parsing.unitName),
		Parameters: []*tree.ClassParameter{},
		SuperTypes: superTypes,
		Children:   nodes,
		Trait: isTrait(nodes),
		Region:     parsing.completeStructure(tree.ClassDeclarationNodeKind),
	}
}

func (parsing *Parsing) parseImplementStatements() (types []tree.TypeName) {
	for token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
	}
	for parsing.isLookingAtKeyword(token.ImplementKeyword) {
		types = append(types, parsing.parseImplementStatement().Trait)
	}
	return types
}

func isTrait(nodes []tree.Node) bool {
	for _, child := range nodes {
		if method, ok := child.(*tree.MethodDeclaration); ok && method.Abstract {
			return true
		}
	}
	return false
}

func (parsing *Parsing) parseTestStatement() tree.Node {
	parsing.beginStructure(tree.TestStatementNodeKind)
	parsing.skipKeyword(token.TestKeyword)
	parsing.skipEndOfStatement()
	statements := parsing.parseStatementBlock()
	return &tree.TestStatement{
		Region:     parsing.completeStructure(tree.TestStatementNodeKind),
		MethodName: parsing.currentMethod.name,
		Body:       statements,
	}
}

func (parsing *Parsing) parseFieldDeclaration() tree.Node {
	parsing.beginStructure(tree.FieldDeclarationNodeKind)
	parsing.skipKeyword(token.HasKeyword)
	name := parsing.parseIdentifier()
	fieldType := parsing.parseTypeName()
	return &tree.FieldDeclaration{
		Region:   parsing.completeStructure(tree.FieldDeclarationNodeKind),
		Name:     name,
		TypeName: fieldType,
	}
}

func (parsing *Parsing) parseLetBinding() *tree.LetBinding {
	parsing.beginStructure(tree.LetBindingNodeKind)
	parsing.skipKeyword(token.LetKeyword)
	names := parsing.parseLetBindingNames()
	parsing.skipOperator(token.AssignOperator)
	value := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.LetBinding{
		Region:     parsing.completeStructure(tree.LetBindingNodeKind),
		Names:      names,
		Expression: value,
	}
}

func (parsing *Parsing) parseLetBindingNames() []*tree.Identifier {
	if token.HasOperatorValue(parsing.token(), token.LeftBracketOperator) {
		return parsing.parseBindingNameList()
	} else {
		return []*tree.Identifier{parsing.parseIdentifier()}
	}
}

func (parsing *Parsing) throwEmptyLetBindingError() {
	parsing.throwError(&diagnostic.RichError{
		CommonReasons: []string{"let binding has no variable names"},
		Error: &diagnostic.UnexpectedTokenError{
			Expected: token.IdentifierTokenName,
			Received: token.RightBracketOperator.String(),
		},
	})
}

func (parsing *Parsing) parseBindingNameList() (names []*tree.Identifier) {
	parsing.skipOperator(token.LeftBracketOperator)
	if parsing.isLookingAtOperator(token.RightBracketOperator) {
		parsing.throwEmptyLetBindingError()
	}
	names = append(names, parsing.parseIdentifier())
	for !token.HasOperatorValue(parsing.token(), token.RightBracketOperator) {
		parsing.skipOperator(token.CommaOperator)
		names = append(names, parsing.parseIdentifier())
	}
	parsing.skipOperator(token.RightBracketOperator)
	return names
}

func (parsing *Parsing) parseLetBindingStatement() tree.Statement {
	parsing.beginStructure(tree.ExpressionStatementNodeKind)
	binding := parsing.parseLetBinding()
	_ = parsing.completeStructure(tree.ExpressionStatementNodeKind)
	return &tree.ExpressionStatement{
		Expression: binding,
	}
}

func (parsing *Parsing) parseImplementStatement() *tree.ImplementStatement {
	parsing.beginStructure(tree.ImplementStatementNodeKind)
	parsing.skipKeyword(token.ImplementKeyword)
	trait := parsing.parseTypeName()
	parsing.skipEndOfStatement()
	return &tree.ImplementStatement{
		Region: parsing.completeStructure(tree.ImplementStatementNodeKind),
		Trait:  trait,
	}
}

func (parsing *Parsing) parseTypeNameList() (names []tree.TypeName) {
	names = append(names, parsing.parseTypeName())
	for parsing.isLookingAtOperator(token.CommaOperator) {
		parsing.skipOperator(token.CommaOperator)
		names = append(names, parsing.parseTypeName())
	}
	return names
}
