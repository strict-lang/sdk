package scope

type compositeScope struct {
	id Id
	children []Scope
}

func Combine(id Id, children ...Scope) Scope {
	return &compositeScope{
		id:       id,
		children: children,
	}
}

func (scope *compositeScope) Id() Id {
	return scope.id
}

func (scope *compositeScope) Lookup(point ReferencePoint) (result EntrySet) {
	for _, scope := range scope.children {
		result = append(result, scope.Lookup(point)...)
	}
	return
}

func (scope compositeScope) Search(filter symbolFilter) (result EntrySet) {
	for _, scope := range scope.children {
		result = append(result, scope.Search(filter)...)
	}
	return
}

func (scope compositeScope) Contains(point ReferencePoint) bool {
	for _, scope := range scope.children {
		if scope.Contains(point) {
			return true
		}
	}
	return false
}
