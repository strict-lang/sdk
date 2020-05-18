package namespace

type Table struct {
	namespaces map[string]Namespace
}

func NewTable() *Table {
	return &Table{
		namespaces: map[string]Namespace{},
	}
}

func (table *Table) Find(qualifiedName string) (Namespace, bool) {
	namespace, ok := table.namespaces[qualifiedName]
	return namespace, ok
}

func (table *Table) FindOrCompute(
	qualifiedName string,
	factory func(string) Namespace) Namespace {

	if namespace, ok := table.namespaces[qualifiedName]; ok {
		return namespace
	}
	createdNamespace := factory(qualifiedName)
	table.namespaces[qualifiedName] = createdNamespace
	return createdNamespace
}

func (table *Table) List() []Namespace {
	var namespaces []Namespace
	for _, namespace := range table.namespaces {
		namespaces = append(namespaces, namespace)
	}
	return namespaces
}
