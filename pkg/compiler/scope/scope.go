package scope

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"


type Id string

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

type Scope interface {
	Id() Id
	Lookup(point ReferencePoint) EntrySet
	Search(filter symbolFilter) EntrySet
	Contains(point ReferencePoint) bool
}

type symbolFactory func() (Symbol, bool)

type MutableScope interface {
	Scope

	Insert(symbol Symbol)
	LookupOrInsert(point ReferencePoint, factory symbolFactory) EntrySet
}

type NamedScope struct {

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


type EmptyScope struct {
	id Id
}

func NewEmptyScope(id Id) Scope {
	return &EmptyScope{id: id}
}

func (scope *EmptyScope) Id() Id {
	return scope.id
}

func (scope *EmptyScope) Insert(symbol Symbol) {}

func (scope *EmptyScope) LookupOrInsert(
	point ReferencePoint, factory symbolFactory) EntrySet {

	return EntrySet{}
}

func (scope *EmptyScope) Lookup(point ReferencePoint) EntrySet {
	return EntrySet{}
}

func (scope *EmptyScope) Contains(point ReferencePoint) bool {
	return false
}

func (scope *EmptyScope) Search(filter symbolFilter) EntrySet {
	return EntrySet{}
}

func LookupClass(scope Scope, point ReferencePoint) (*Class, bool) {
	if entries := scope.Lookup(point); !entries.IsEmpty() {
		first := entries.First().Symbol
		class, isClass := first.(*Class)
		return class, isClass
	}
	return nil, false
}