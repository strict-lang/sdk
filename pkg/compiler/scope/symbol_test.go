package scope

import "github.com/strict-lang/sdk/pkg/compiler/input"

type testSymbol struct {
	name              string
	declarationOffset input.Offset
}

func (symbol *testSymbol) Name() string {
	return symbol.name
}

func (symbol *testSymbol) String() string {
	return symbol.name
}

func (symbol *testSymbol) DeclarationOffset() input.Offset {
	return symbol.declarationOffset
}

func NewTestSymbol(name string) Symbol {
	return &testSymbol{name: name}
}

func NewPositionedTestSymbol(name string, offset input.Offset) Symbol {
	return &testSymbol{
		name:              name,
		declarationOffset: offset,
	}
}
