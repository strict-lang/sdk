package code

import (
	"errors"
	"fmt"
)

var (
	ErrSymbolExists = errors.New("constantpool already exists")
	ErrNoSuchSymbol = errors.New("constantpool does not exist")
)

const (
	ScopedTypeAttribute AttributeKind = iota
	ScopedFieldAttribute
	ScopedMethodAttribute
)

type Scope struct {
	parent     *Scope
	name       string
	depth      int
	childCount int
	symbols    map[string]Attribute
}

func NewRootScope() *Scope {
	return &Scope{
		parent:     nil,
		name:       "@",
		depth:      0,
		childCount: 0,
		symbols:    make(map[string]Attribute),
	}
}

func (scope *Scope) PutSymbol(symbol string, attribute Attribute) error {
	if scope.ContainsSymbol(symbol) {
		return ErrSymbolExists
	}
	scope.symbols[symbol] = attribute
	return nil
}

func (scope *Scope) RemoveSymbol(symbol string) (Attribute, error) {
	if scoped, ok := scope.symbols[symbol]; ok {
		delete(scope.symbols, symbol)
		return scoped, nil
	}
	if scope.parent == nil {
		return Attribute{}, ErrNoSuchSymbol
	}
	return scope.parent.RemoveSymbol(symbol)
}

func (scope *Scope) LookupSymbol(symbol string) (Attribute, bool) {
	if scoped, ok := scope.symbols[symbol]; ok {
		return scoped, true
	}
	if scope.parent == nil {
		return Attribute{}, false
	}
	return scope.parent.LookupSymbol(symbol)
}

func (scope *Scope) ContainsSymbol(symbol string) bool {
	if _, ok := scope.symbols[symbol]; ok {
		return true
	}
	if scope.parent == nil {
		return false
	}
	return scope.parent.ContainsSymbol(symbol)
}

func (scope *Scope) Name() string {
	return scope.name
}

func (scope *Scope) IsRoot() bool {
	return scope.parent == nil
}

func (scope *Scope) Parent() *Scope {
	return scope.parent
}

func (scope *Scope) NewNamedChild(name string) *Scope {
	var childName string
	if scope.IsRoot() {
		childName = fmt.Sprintf("@%s", name)
	} else {
		childName = fmt.Sprintf("%s/%s", scope.name, name)
	}
	scope.childCount++
	return &Scope{
		parent:  scope,
		name:    childName,
		depth:   scope.depth + 1,
		symbols: make(map[string]Attribute),
	}
}

func (scope *Scope) NewChild() *Scope {
	return scope.NewNamedChild(fmt.Sprintf("child-%d", scope.childCount+1))
}