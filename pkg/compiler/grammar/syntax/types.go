package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

func (parsing *Parsing) couldBeLookingAtTypeName() bool {
	if !token.IsIdentifierToken(parsing.token()) {
		return false
	}
	peek := parsing.peek()
	if token.IsIdentifierToken(peek) {
		return true
	}
	return token.HasOperatorValue(peek, token.SmallerOperator) ||
		token.HasOperatorValue(peek, token.LeftBracketOperator)
}

// parseTypeName is a recursive method that parses type names. When calling
// this method, the types primary name is the value of the 'last' token.
func (parsing *Parsing) parseTypeName() tree.TypeName {
	parsing.beginStructure(tree.TypeNameNodeGroup)
	base := parsing.parseIdentifier()
	return parsing.parseTypeNameFromBaseIdentifier(base.Value)
}

func (parsing *Parsing) parseTypeNameFromBaseIdentifier(base string) tree.TypeName {
	typeName := parsing.parseIncompleteGenericOrConcreteType(base)
	if token.HasOperatorValue(parsing.token(), token.LeftBracketOperator) {
		return parsing.parseListTypeName(typeName)
	}
	if token.HasOperatorValue(parsing.token(), token.QuestionMarkOperator) {
		return parsing.parseOptionalTypeName(typeName)
	}
	parsing.completeStructure(tree.WildcardNodeKind)
	return typeName
}

func (parsing *Parsing) expectBaseName(name token.Token) {
	if !token.IsIdentifierToken(name) {
		parsing.throwError(&diagnostic.RichError{
			Error: &diagnostic.UnexpectedTokenError{
				Expected: "name of type",
				Received: name.Value(),
			},
		})
	}
}

func (parsing *Parsing) parseOptionalTypeName(base tree.TypeName) tree.TypeName {
	return &tree.OptionalTypeName{
		Region:   parsing.createRegionOfCurrentStructure(),
		TypeName: base,
	}
}

func (parsing *Parsing) parseIncompleteGenericOrConcreteType(base string) tree.TypeName {
	if token.HasOperatorValue(parsing.token(), token.SmallerOperator) {
		return parsing.parseIncompleteGenericTypeName(base)
	}
	if token.HasOperatorValue(parsing.token(), token.QuestionMarkOperator) {
		return parsing.parseIncompleteOptionalTypeName(base)
	}
	return parsing.parseIncompleteConcreteTypeName(base)
}

func (parsing *Parsing) parseIncompleteOptionalTypeName(base string) tree.TypeName {
	parsing.updateTopStructureKind(tree.OptionalTypeNameNodeKind)
	return &tree.OptionalTypeName{
		Region: parsing.createRegionOfCurrentStructure(),
		TypeName: &tree.ConcreteTypeName{
			Name:   base,
			Region: parsing.createRegionOfCurrentStructure(),
		},
	}
}

func (parsing *Parsing) parseIncompleteConcreteTypeName(base string) tree.TypeName {
	parsing.updateTopStructureKind(tree.ConcreteTypeNameNodeKind)
	return &tree.ConcreteTypeName{
		Name:   base,
		Region: parsing.createRegionOfCurrentStructure(),
	}
}

func (parsing *Parsing) parseIncompleteGenericTypeName(base string) tree.TypeName {
	parsing.updateTopStructureKind(tree.GenericTypeNameNodeKind)
	parsing.skipOperator(token.SmallerOperator)
	generic := parsing.parseTypeName()
	parsing.skipEndOfGenericTypeName()
	return &tree.GenericTypeName{
		Name:    base,
		Generic: generic,
		Region:  parsing.createRegionOfCurrentStructure(),
	}
}

func (parsing *Parsing) skipEndOfGenericTypeName() {
	parsing.skipOperator(token.GreaterOperator)
}

func (parsing *Parsing) parseListTypeName(base tree.TypeName) tree.TypeName {
	parsing.skipOperator(token.LeftBracketOperator)
	parsing.skipOperator(token.RightBracketOperator)
	if token.HasOperatorValue(parsing.token(), token.LeftBracketOperator) {
		beginOffset := parsing.peekStructure().beginOffset
		return parsing.parseListTypeName(&tree.ListTypeName{
			Element: base,
			Region:  input.CreateRegion(beginOffset, parsing.offset()),
		})
	}
	return &tree.ListTypeName{
		Element: base,
		Region:  parsing.completeStructure(tree.WildcardNodeKind),
	}
}
