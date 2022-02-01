package dns

import "errors"

const (
	// TypeA An ipv4 host address
	TypeA Type = 1
	// TypeNS An authoritative name server
	TypeNS Type = 2
	// TypeMD a mail destination (Obsolete - use MX)
	TypeMD Type = 3
	// TypeMF a mail forwarder (Obsolete - use MX)
	TypeMF Type = 4
	// TypeCNAME the canonical name for an alias
	TypeCNAME Type = 5
	// TypeSOA marks the start of a zone of authority
	TypeSOA Type = 6
	// TypeMB a mailbox domain name (EXPERIMENTAL)
	TypeMB Type = 7
	// TypeMG a mail group member (EXPERIMENTAL)
	TypeMG Type = 8
	// TypeMR a mail rename domain name (EXPERIMENTAL)
	TypeMR Type = 9
	// TypeNULL a null RR (EXPERIMENTAL)
	TypeNULL Type = 10
	// TypeWKS a well known service description
	TypeWKS Type = 11
	// TypePTR a domain name pointer
	TypePTR Type = 12
	// TypeHINFO host information
	TypeHINFO Type = 13
	// TypeMINFO mailbox or mail list information
	TypeMINFO Type = 14
	// TypeMX mail exchange
	TypeMX Type = 15
	// TypeTXT text strings
	TypeTXT Type = 16
	// TypeAAAA An ipv6 host address
	TypeAAAA Type = 28

	// QTypes

	// QTypeAXFR A request for a transfer of an entire zone
	QTypeAXFR Type = 252
	// QTypeMAILB A request for mailbox-related records (MB, MG or MR)
	QTypeMAILB Type = 253
	// QTypeMAILA A request for mail agent RRs (Obsolete - see MX)
	QTypeMAILA Type = 254
	// QTypeAll A request for all records
	QTypeAll Type = 255

	// ClassIN the internet
	ClassIN Class = 1
	// ClassCS the CSNET class (Obsolete - used only for examples in some obsolete RFCs)
	ClassCS Class = 2
	// ClassChaos the CHAOS class
	ClassChaos Class = 3
	// ClassHS Hesiod
	ClassHS Class = 4

	// QClasses

	// QClassAny any class
	QClassAny Class = 255
)

type Type uint16

func (t Type) encode(w writeOffsetter, c *compressionCache) error {
	return writeUint16(w, uint16(t))
}

func (t *Type) decode(r readSeekOffsetter) error {
	n, err := readUint16(r)
	if err != nil {
		return err
	}
	*t = Type(n)
	return nil
}

type Class uint16

func (cl Class) encode(w writeOffsetter, c *compressionCache) error {
	return writeUint16(w, uint16(cl))
}

func (cl *Class) decode(r readSeekOffsetter) error {
	n, err := readUint16(r)
	if err != nil {
		return err
	}
	*cl = Class(n)
	return nil
}

type ResourceRecord struct {
	// an owner name, i.e., the name of the node to which
	// this resource record pertains.
	Name Name
	// two octets containing one of the RR TYPE codes
	Type Type
	// two octets containing one of the RR CLASS codes
	Class Class

	TTL uint32

	Data interface{}
}

func (rr *ResourceRecord) encode(w writeOffsetter, c *compressionCache) error {
	if err := rr.Name.encode(w, c); err != nil {
		return err
	}

	if err := rr.Type.encode(w, c); err != nil {
		return err
	}

	if err := rr.Class.encode(w, c); err != nil {
		return err
	}

	if err := writeUint32(w, rr.TTL); err != nil {
		return err
	}

	dEnc, ok := rr.Data.(encoder)
	if !ok {
		return errors.New("data does not implement encoder interface")
	}

	return dEnc.encode(w, c)
}

func (rr *ResourceRecord) decode(r readSeekOffsetter) error {

	if err := rr.Name.decode(r); err != nil {
		return err
	}

	if err := rr.Type.decode(r); err != nil {
		return err
	}

	if err := rr.Class.decode(r); err != nil {
		return err
	}

	ttl, err := readUint32(r)
	if err != nil {
		return err
	}
	rr.TTL = ttl

	data := getRecordData(rr.Type)
	dec, ok := data.(decoder)
	if !ok {
		return errors.New("invalid record data")
	}

	if err := dec.decode(r); err != nil {
		return err
	}

	rr.Data = dec

	return nil
}

func getRecordData(t Type) interface{} {
	switch t {
	case TypeA:
		return ARecordData{}
	case TypeNS:
		return NSRecordData{}
	default:
		return unsupportedRecordData{}
	}

}

type nameRecordData struct {
	Name Name
}

func (n *nameRecordData) encode(w writeOffsetter, c *compressionCache) error {
	buf := newOffsetWriter(w.Offset())
	if err := n.Name.encode(buf, c); err != nil {
		return err
	}
	if err := writeUint16(w, buf.Len()); err != nil {
		return err
	}

	_, err := buf.WriteTo(w)

	return err
}

func (n *nameRecordData) decode(r readSeekOffsetter) error {
	_, err := readUint16(r)
	if err != nil {
		return err
	}
	return n.Name.decode(r)
}

func (n nameRecordData) String() string {
	return n.Name.String()
}

type NSRecordData struct {
	nameRecordData
}

type CNameRecordData struct {
	nameRecordData
}

type PTRRecordData struct {
	nameRecordData
}
