package scope

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// ShadowingPolicy determines whether a LocalScope can shadow its parents
// definitions. It has to states and may be represented by a boolean.
type ShadowingPolicy bool

const (
	AllowShadowing ShadowingPolicy = true
	ForbidShadowing ShadowingPolicy = false
)

// LocalScope is a scope that is within a function or block of code. It can
// shadow definitions in the parent scopes and is position aware, meaning that
// the scope's entries are only visible when the ReferencePoint's location is
// after the entry's location. Non-LocalScope's act differently and that's fine
// since they do not contain code with semantics bound to the program order.
//
// LocalScopes come in to flavours, shadowing and non-shadowing. The flavour is
// configured when creating a new instance. Shadowing scopes can shadow
// definitions in the parents while non-shadowing ones can't.
//
// A current problem is, that scopes do not care what kind of scope their parent is,
// if the parent is a shadowing LocalScope, that shadows variable x' and the child
// also shadows variable x', then a method would contain x' two times.
// In practice, we only use shadowing scopes when creating the method scope, in order
// to have parameters that share a name with a classes attribute.
type LocalScope struct {
	id              Id
	parent          Scope
	region          input.Region
	entries         map[string]Entry
	shadowingPolicy ShadowingPolicy
}

// NewShadowingLocalScope creates a LocalScope that may shadow definitions.
func NewShadowingLocalScope(
	id Id,
	region input.Region,
	parent Scope) MutableScope {

	return &LocalScope{
		id:              id,
		parent:          parent,
		region:          region,
		entries:         map[string]Entry{},
		shadowingPolicy: AllowShadowing,
	}
}

// NewLocalScope creates a non-shadowing LocalScope.
func NewLocalScope(
	id Id,
	region input.Region,
	parent Scope) MutableScope {

	return &LocalScope{
		id:              id,
		parent:          parent,
		region:          region,
		entries:         map[string]Entry{},
		shadowingPolicy: ForbidShadowing,
	}
}

func createReferencePointFromSymbol(symbol Symbol) ReferencePoint {
	return ReferencePoint{
		name:           symbol.Name(),
		ignorePosition: true,
	}
}

func (scope *LocalScope) Id() Id {
	return scope.id
}

// Inserts a symbol into the local scope, if the symbol has not yet been
// defined or may shadow a previous definition.
func (scope *LocalScope) Insert(symbol Symbol) {
	point := createReferencePointFromSymbol(symbol)
	if scope.parent.Contains(point) {
		scope.maybeShadow(point, symbol)
	} else {
		scope.insertToLocal(symbol)
	}
}

func (scope *LocalScope) LookupOrInsert(
	point ReferencePoint,
	factory symbolFactory) EntrySet {

	if found := scope.Lookup(point); len(found) != 0 {
		return found
	}
	if symbol, ok := factory(); ok {
		createdEntry := scope.insertToLocal(symbol)
		return EntrySet{createdEntry }
	}
	return EntrySet{}
}

// Lookup searches for all symbols that match the ReferencePoint's name and can be
// accessed from the point's position (given that the position is not ignored).
func (scope *LocalScope) Lookup(point ReferencePoint) EntrySet {
	parentEntries := scope.parent.Lookup(point)
	if len(parentEntries) != 0 && scope.canShadow() {
		ownEntries := scope.lookupOwn(point)
		return append(parentEntries, ownEntries...)
	}
	return scope.lookupOwn(point)
}

func (scope *LocalScope) lookupOwn(point ReferencePoint) EntrySet {
	symbol := scope.entries[point.name]
	if point.ignorePosition || canSeeEntry(point, symbol) {
		return EntrySet{symbol}
	}
	return EntrySet{}
}

// Shadows a symbol that has been found in the parent scope, if the current
// scope does not already contain it and shadowing is allowed.
func (scope *LocalScope) maybeShadow(point ReferencePoint, symbol Symbol) {
	if !scope.containsDirectly(point) && scope.canShadow() {
		scope.insertToLocal(symbol)
	}
}

func (scope *LocalScope) canShadow() bool {
	return scope.shadowingPolicy == AllowShadowing
}

func (scope* LocalScope) createEntry(symbol Symbol) Entry {
	return Entry{
		symbol:   symbol,
		position: symbol.DeclarationOffset(),
		scopeId:  scope.id,
	}
}

func (scope* LocalScope) insertToLocal(symbol Symbol) Entry {
	entry := scope.createEntry(symbol)
	scope.entries[symbol.Name()] = entry
	return entry
}

func (scope *LocalScope) containsDirectly(point ReferencePoint) bool {
	_, ok := scope.entries[point.name]
	return ok
}

// Returns whether the local scope has a symbol that both matches the name
// of the reference point and has been declared before the point. This is
// done when entities are looked up who's location and order is of semantic
// relevance. Examples are variables: If a variable x' accesses y' and both
// are declared in the same scope, then y' has to be declared before x'.
func (scope *LocalScope) containsBeforePoint(point ReferencePoint) bool {
	symbol, ok := scope.entries[point.name]
	if !ok {
		return false
	}
	return canSeeEntry(point, symbol)
}

func canSeeEntry(point ReferencePoint, entry Entry) bool {
	return entry.position < point.position
}

// Contains returns whether the scope contains any symbols matching
// the point's name, that can also be accessed from the points position
// (given that the position is not ignored).
func (scope *LocalScope) Contains(point ReferencePoint) bool {
	if scope.parent.Contains(point) {
		return true
	}
	if point.ignorePosition {
		return scope.containsDirectly(point)
	}
	return scope.containsBeforePoint(point)
}

// Search traverses the scope and its parent scopes and applies the filter on
// every entry. The result is a set of entries that have been approved by the
// filter. The positional visibility of entries is not taken into account.
func (scope *LocalScope) Search(filter symbolFilter) EntrySet {
	own := scope.searchOwn(filter)
	parents := scope.parent.Search(filter)
	return append(parents, own...)
}

func (scope *LocalScope) searchOwn(filter symbolFilter) (result EntrySet) {
	for _, entry := range scope.entries {
		if filter(entry.symbol) {
			result = append(result, entry)
		}
	}
	return result
}
