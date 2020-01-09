package scope

type EmptyScope struct {
	id Id
}

func NewEmptyScope(id Id) Scope {
	return &EmptyScope{id: id}
}

func (scope *EmptyScope) Id() Id {
	return scope.id
}

func (scope *EmptyScope) Insert(Symbol) {}

func (scope *EmptyScope) LookupOrInsert(ReferencePoint, symbolFactory) EntrySet {
	return EntrySet{}
}

func (scope *EmptyScope) Lookup(ReferencePoint) EntrySet {
	return EntrySet{}
}

func (scope *EmptyScope) Contains(ReferencePoint) bool {
	return false
}

func (scope *EmptyScope) Search(symbolFilter) EntrySet {
	return EntrySet{}
}
