package resolver

import (
	"net"

	"github.com/davidseybold/dns-resolver/dns"
	"github.com/davidseybold/dns-resolver/resolver/cache"
)

type Resolver struct {
	cache           *cache.Cache
	sBelt           srvList
	pendingRequests map[string]*request
}

func NewResolver() *Resolver {
	return &Resolver{
		cache: cache.New(),
	}
}

func (r *Resolver) LookupHost(name string) ([]net.IP, error) {
	records, err := r.Lookup(dns.DomainName(name), dns.ClassIN, dns.TypeA)
	if err != nil {
		return []net.IP{}, err
	}
	ips := []net.IP{}
	for _, record := range records {
		aRecord := dns.ParseARecord(record.RData, record.RDLength)
		ips = append(ips, aRecord.Address)
	}
	return ips, nil
}

func (r *Resolver) LookupAddress(addr net.IP) ([][]byte, error) {
	addrCopy := make([]byte, len(addr))
	copy(addrCopy, addr)
	for i, j := 0, len(addrCopy)-1; i < j; i, j = i+1, j-1 {
		addrCopy[i], addrCopy[j] = addrCopy[j], addrCopy[i]
	}
	name := append(addrCopy, []byte(".in-addr.arpa")...)

	records, err := r.Lookup(name, dns.ClassIN, dns.TypePTR)
	if err != nil {
		return [][]byte{}, err
	}

	names := [][]byte{}
	for _, record := range records {
		ptrRecord := dns.ParsePtrRecord(record.RData, record.RDLength)
		names = append(names, ptrRecord.DomainName)
	}

	return names, nil
}

func (r *Resolver) Lookup(qName dns.DomainName, qClass dns.Class, qType dns.Type) ([]dns.ResourceRecord, error) {
	return r.lookup(dns.Question{
		DomainName: qName,
		Class:      qClass,
		Type:       qType,
	})
}

func (r *Resolver) lookup(question dns.Question) ([]dns.ResourceRecord, error) {

	return nil, nil
}
