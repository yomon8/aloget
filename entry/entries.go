package entry

type Entries []*Entry

func (es Entries) Len() int {
	return len(es)
}

func (es Entries) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es Entries) Less(i, j int) bool {
	return es[i].RequestTime.Before(es[j].RequestTime)
}

func (es Entries) GetAllEntries() []*Entry {
	entries := make([]*Entry, es.Len())
	for i, e := range es {
		entries[i] = e
	}
	return entries
}
