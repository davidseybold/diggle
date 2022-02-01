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
	records, err := r.Lookup(dns.NewName(name), dns.ClassIN, dns.TypeA)
	if err != nil {
		return []net.IP{}, err
	}
	ips := []net.IP{}
	for _, r := range records {
		a, ok := r.Data.(dns.ARecordData)
		if ok {
			ips = append(ips, a.Address)
		}
	}
	return ips, nil
}

func (r *Resolver) LookupAddress(addr net.IP) ([][]byte, error) {

	return nil, nil
}

func (r *Resolver) Lookup(qName dns.Name, qClass dns.Class, qType dns.Type) ([]dns.ResourceRecord, error) {
	return r.lookup(dns.Question{
		Name:  qName,
		Class: qClass,
		Type:  qType,
	})
}

func (r *Resolver) lookup(question dns.Question) ([]dns.ResourceRecord, error) {

	return nil, nil
}
