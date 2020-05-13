package scope

type Cache struct {
	entries map[string] cacheEntry
}

type cacheEntry struct {
	scope Scope
}
