package isolate

type ThreadLocalPropertyTable struct {
	properties map[string] interface{}
}

func (table *ThreadLocalPropertyTable) Insert(name string, value interface{}) {
	table.properties[name] = value
}

func (table *ThreadLocalPropertyTable) Lookup(name string) (interface{}, bool) {
	value, ok := table.properties[name]
	return value, ok
}
