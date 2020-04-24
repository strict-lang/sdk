package diagnostic

type Diagnostics struct {
	entries map[string][]Entry
}

func (diagnostics *Diagnostics) ListEntries() (values []Entry) {
	for _, entries := range diagnostics.entries {
		for _, entry := range entries {
			values = append(values, entry)
		}
	}
	return values
}
