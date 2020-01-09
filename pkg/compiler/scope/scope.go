package scope

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type Id string

type Scope interface {
	Id() Id
	Lookup(point ReferencePoint) EntrySet
	Search(filter symbolFilter) EntrySet
	Contains(point ReferencePoint) bool
}

type ReferencePoint struct {
	name           string
	position       input.Offset
	ignorePosition bool
}

type Entry struct {
	Symbol   Symbol
	position input.Offset
	scopeId  Id
}

type EntrySet []Entry

func (set EntrySet) First() Entry {
	return set[0]
}

func (set EntrySet) IsEmpty() bool {
	return len(set) == 0
}

type symbolFilter func(Symbol) bool

type symbolFactory func() (Symbol, bool)

type MutableScope interface {
	Scope

	Insert(symbol Symbol)
	LookupOrInsert(point ReferencePoint, factory symbolFactory) EntrySet
}

func NewReferencePoint(name string) ReferencePoint {
	return ReferencePoint{
		name:           name,
		ignorePosition: true,
	}
}

func NewReferencePointWithPosition(
	name string, position input.Offset) ReferencePoint {

	return ReferencePoint{
		name:           name,
		position:       position,
		ignorePosition: false,
	}
}

func LookupClass(scope Scope, point ReferencePoint) (*Class, bool) {
	if entries := scope.Lookup(point); !entries.IsEmpty() {
		first := entries.First().Symbol
		class, isClass := first.(*Class)
		return class, isClass
	}
	return nil, false
}