package resolver

import (
	"net"
	"time"

	"github.com/davidseybold/dns-resolver/dns"
)

type srv struct {
	Priority int
	Used     bool
}

type addrScore struct {
	BattingAvg         float32
	MedianResponseTime time.Duration
}

type sList struct {
	ZoneName   dns.Name
	ZoneNS     []srv
	NSAddr     map[string]net.IP
	AddrScores map[string]addrScore
}

// func (n srvList) SortByPriority() {
// 	sort.Slice(n, func(i, j int) bool {
// 		return n[i].Priority < n[j].Priority
// 	})
// }
