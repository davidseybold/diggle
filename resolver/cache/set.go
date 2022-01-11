package cache

type recordSet map[string]cacheRecord

func (r recordSet) Add(records ...cacheRecord) {
	for i := range records {
		key := records[i].Hash()
		r[key] = records[i]
	}
}

func (r recordSet) Delete(rec cacheRecord) bool {
	key := rec.Hash()
	delete(r, key)
	return true
}

func (r recordSet) Records() []cacheRecord {
	records := []cacheRecord{}
	for _, record := range r {
		records = append(records, record)
	}
	return records
}
