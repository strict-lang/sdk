package diagnostic

type Diagnostics struct {
	entries []Entry
}

func Empty() *Diagnostics {
	return &Diagnostics{}
}

func (diagnostics *Diagnostics) Merge(target *Diagnostics) *Diagnostics {
	mergedEntries := append(diagnostics.entries, target.entries...)
	return &Diagnostics{entries: mergedEntries}
}

func (diagnostics *Diagnostics) ListEntries() (values []Entry) {
	copied := make([]Entry, len(diagnostics.entries))
	copy(copied, values)
	return copied
}

func Merge(target ...*Diagnostics) *Diagnostics {
	var identity []Entry
	for _, diagnostic := range target {
		identity = append(identity, diagnostic.entries...)
	}
	return &Diagnostics{entries: identity}
}
