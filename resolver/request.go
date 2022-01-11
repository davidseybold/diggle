package resolver

import (
	"net"
	"time"

	"github.com/davidseybold/dns-resolver/dns"
)

type request struct {
	ID          int16
	StartTime   time.Time
	StepCounter int

	SName       dns.DomainName
	SType       dns.Type
	SClass      dns.Class
	NSList      []string
	SBelt       []string
	NSAddresses map[string]net.IP
}

func (r *request) Start() error {

}
