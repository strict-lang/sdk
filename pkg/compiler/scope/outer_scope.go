package scope

type OuterScope struct {
	id Id
	parent Scope
	entries map[string] Entry
}

func NewOuterScope(id Id, parent Scope) MutableScope {
	return &OuterScope{
		id:      parent.Id().NewChildId(id),
		parent:  parent,
		entries: map[string] Entry{},
	}
}

func (scope *OuterScope) Id() Id {
	return scope.id
}

// Inserts a Symbol into the local scope, if the Symbol has not
// isn't defined. This call does never shadow symbols.
func (scope *OuterScope) Insert(symbol Symbol) {
	point := createReferencePointFromSymbol(symbol)
	if !scope.parent.Contains(point) {
		scope.insertToSelf(symbol)
	}
}

func (scope *OuterScope) LookupOrInsert(
	point ReferencePoint,
	factory symbolFactory) EntrySet {

	if found := scope.Lookup(point); len(found) != 0 {
		return found
	}
	if symbol, ok := factory(); ok {
		createdEntry := scope.insertToSelf(symbol)
		return EntrySet{createdEntry }
	}
	return EntrySet{}
}

// Lookup searches for all symbols that match the points name.
func (scope *OuterScope) Lookup(point ReferencePoint) EntrySet {
	if entries := scope.parent.Lookup(point); !entries.IsEmpty() {
		return entries
	}
	return scope.lookupOwn(point)
}

func (scope *OuterScope) lookupOwn(point ReferencePoint) EntrySet {
	if symbol, ok := scope.entries[point.name]; ok {
		return EntrySet{symbol}
	}
	return EntrySet{}
}

func (scope* OuterScope) createEntry(symbol Symbol) Entry {
	return Entry{
		Symbol:   symbol,
		position: symbol.DeclarationOffset(),
		scopeId:  scope.id,
	}
}

func (scope* OuterScope) insertToSelf(symbol Symbol) Entry {
	entry := scope.createEntry(symbol)
	scope.entries[symbol.Name()] = entry
	return entry
}

func (scope *OuterScope) containsDirectly(point ReferencePoint) bool {
	_, ok := scope.entries[point.name]
	return ok
}

// Contains returns whether the scope contains any symbols matching
// the point's DeclarationName, that can also be accessed from the points position
// (given that the position is not ignored).
func (scope *OuterScope) Contains(point ReferencePoint) bool {
	if scope.parent.Contains(point) {
		return true
	}
	return scope.containsDirectly(point)
}

// Search traverses the scope and its parent scopes and applies the filter on
// every entry. The result is a set of entries that have been approved by the
// filter.
func (scope *OuterScope) Search(filter symbolFilter) EntrySet {
	own := scope.searchOwn(filter)
	parents := scope.parent.Search(filter)
	return append(parents, own...)
}

func (scope *OuterScope) searchOwn(filter symbolFilter) (result EntrySet) {
	for _, entry := range scope.entries {
		if filter(entry.Symbol) {
			result = append(result, entry)
		}
	}
	return result
}
