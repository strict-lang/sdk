package scope

import "sync"

type NamespaceTable struct {
	entries map[string] cacheEntry
	mutex sync.RWMutex
}

var globalTable = &NamespaceTable{
	entries: map[string]cacheEntry{},
	mutex: sync.RWMutex{},
}

func GlobalNamespaceTable() *NamespaceTable {
	return globalTable
}

type cacheEntry struct {
	namespace *Namespace
}

func (cache *NamespaceTable) Insert(name string, namespace *Namespace) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.entries[name] = cacheEntry{namespace: namespace}
}

func (cache *NamespaceTable) Lookup(name string) (*Namespace, bool) {
	cache.mutex.RLock()
	cache.mutex.RUnlock()
	if entry, ok := cache.entries[name]; ok {
		return entry.namespace, true
	}
	return nil, false
}