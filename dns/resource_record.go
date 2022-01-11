package dns

import (
	"net"
)

type ResourceRecord struct {
	// an owner name, i.e., the name of the node to which
	// this resource record pertains.
	Name Name
	// two octets containing one of the RR TYPE codes
	Type Type
	// two octets containing one of the RR CLASS codes
	Class Class

	TTL uint32
	// specifies the length in octets of the RDATA field
	rdLength uint16
	// a variable length string of octets that describes the
	// resource.  The format of this information varies
	// according to the TYPE and CLASS of the resource record
	rData []byte

	nameLookup nameLookup
}

func (rr *ResourceRecord) marshal(w writer) error {
	err := w.WriteName(rr.Name)
	if err != nil {
		return err
	}

	err = w.WriteData(uint16(rr.Type))
	if err != nil {
		return err
	}

	err = w.WriteData(uint16(rr.Class))
	if err != nil {
		return err
	}

	err = w.WriteData(rr.TTL)
	if err != nil {
		return err
	}

	err = w.WriteData(rr.rdLength)
	if err != nil {
		return err
	}

	err = w.WriteData(rr.rData)
	if err != nil {
		return err
	}

	return nil
}

func (rr *ResourceRecord) unmarshal(r reader) error {
	name, err := r.ReadName()
	if err != nil {
		return err
	}

	rType, err := r.ReadUint16()
	if err != nil {
		return err
	}

	rClass, err := r.ReadUint16()
	if err != nil {
		return err
	}

	ttl, err := r.ReadUint32()
	if err != nil {
		return err
	}

	rdLength, err := r.ReadUint16()
	if err != nil {
		return err
	}

	rData, err := r.ReadNBytes(int(rdLength))
	if err != nil {
		return nil
	}

	rr = &ResourceRecord{
		Name:     name,
		Type:     Type(rType),
		Class:    Class(rClass),
		TTL:      ttl,
		rdLength: rdLength,
		rData:    rData,
	}
	return nil
}

// An A record
type ARecord struct {
	ResourceRecord
	Address net.IP
}

func (a *ARecord) Decode() error {
	a.Address = net.IP(a.rData)
	return nil
}

func (a *ARecord) Encode() error {
	a.rData = a.Address
	a.rdLength = uint16(len(a.rData))
	return nil
}

// A Ptr record
type PtrRecord struct {
	ResourceRecord
	DomainName Name
}

// A NS record
type NSRecord struct {
	ResourceRecord
	DomainName Name
}

type CNameRecord struct {
	ResourceRecord
	DomainName Name
}

type MXRecord struct {
	ResourceRecord
	DomainName Name
	Preference uint16
}

type SOARecord struct {
	ResourceRecord
	MName   Name
	RName   Name
	Serial  uint32
	Refresh int32
	Retry   int32
	Expire  int32
	Minimum uint32
}
