package isolate

import "log"

type ThreadLocalPropertyTable struct {
	properties map[string]interface{}
}

func (table *ThreadLocalPropertyTable) Insert(name string, value interface{}) {
	table.properties[name] = value
}

func (table *ThreadLocalPropertyTable) Lookup(name string) (interface{}, bool) {
	value, ok := table.properties[name]
	return value, ok
}

func (table *ThreadLocalPropertyTable) Log() {
	log.Printf("properties:")
	for key, property := range table.properties {
		log.Printf("- %s: %v", key, property)
	}
}

func NewThreadLocalPropertyTable() *ThreadLocalPropertyTable {
	return &ThreadLocalPropertyTable{
		properties: map[string]interface{}{},
	}
}
