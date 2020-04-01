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

func (diagnostics *Diagnostics) PrintEntries(printer Printer) {
	for _, entries := range diagnostics.entries {
		for _, entry := range entries {
			entry.PrintColored(printer)
		}
		printer.PrintLine("")
	}
}
