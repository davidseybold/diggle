package cache

import (
	"fmt"
	"sync"
	"time"

	hashstructure "github.com/mitchellh/hashstructure/v2"

	"github.com/davidseybold/dns-resolver/dns"
)

type cacheRecord struct {
	ExpirationTime time.Time
	ResourceRecord dns.ResourceRecord
}

func (c cacheRecord) Hash() string {
	hash, _ := hashstructure.Hash(c.ResourceRecord, hashstructure.FormatV2, nil)
	return fmt.Sprint(hash)
}

func (c cacheRecord) IsExpired() bool {
	return c.ExpirationTime.After(time.Now())
}

type Cache struct {
	mu *sync.RWMutex
	c  map[string]recordSet
}

func New() *Cache {
	return &Cache{
		mu: &sync.RWMutex{},
		c:  make(map[string]recordSet),
	}
}

func (r *Cache) Get(name string) ([]dns.ResourceRecord, bool) {
	cacheRecords, ok := r.get(name)
	if !ok {
		return []dns.ResourceRecord{}, false
	}

	dnsRecords := []dns.ResourceRecord{}
	for _, cRec := range cacheRecords {
		dnsRecords = append(dnsRecords, cRec.ResourceRecord)
	}

	return dnsRecords, len(dnsRecords) > 0
}

func (r *Cache) get(name string) ([]cacheRecord, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	set, ok := r.c[name]
	if !ok {
		return []cacheRecord{}, false
	}
	cacheRecords := set.Records()

	validRecords := []cacheRecord{}
	for i := range cacheRecords {
		if cacheRecords[i].IsExpired() {
			set.Delete(cacheRecords[i])
		} else {
			validRecords = append(validRecords, cacheRecords[i])
		}
	}

	return validRecords, len(validRecords) > 0
}

func (r *Cache) Add(name string, records ...dns.ResourceRecord) {
	timeIn := time.Now()
	cacheRecords := []cacheRecord{}
	for i := range records {
		cacheRecords = append(cacheRecords, cacheRecord{
			ExpirationTime: timeIn.Add(time.Duration(records[i].TTL) * time.Second),
			ResourceRecord: records[i],
		})
	}
	r.add(name, cacheRecords)
}

func (r *Cache) add(name string, cr []cacheRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.c[name]
	if !ok {
		r.c[name] = make(recordSet)
	}
	r.c[name].Add(cr...)
}

func (r *Cache) Query(name string, qType dns.Type, qclass dns.Class) ([]dns.ResourceRecord, bool) {
	records, ok := r.Get(name)
	if !ok {
		return []dns.ResourceRecord{}, false
	}
	filteredRecords := []dns.ResourceRecord{}
	for i := range records {
		if records[i].Class == qclass && records[i].Type == qType {
			filteredRecords = append(filteredRecords, records[i])
		}
	}
	return filteredRecords, len(filteredRecords) > 0
}
