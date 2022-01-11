package dns

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
)

const (
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
type Class uint16

type Packet struct {
	Header      Header
	Questions   []Question
	Answers     []ResourceRecord
	Authorities []ResourceRecord
	Additional  []ResourceRecord
}

type marshaller interface {
	marshal(writer) error
}

type unmarshaller interface {
	unmarshal(reader) error
}

func (p *Packet) Unmarshal(buf []byte) error {
	r := newPacketReader(buf)

	var h *Header
	err := h.unmarshal(r)
	if err != nil {
		return err
	}

	questions := []Question{}
	for i := 0; i < int(h.QDCount); i++ {
		var q *Question
		err = q.unmarshal(r)
		if err != nil {
			return err
		}
		questions = append(questions, *q)
	}

	answers := []ResourceRecord{}
	for i := 0; i < int(h.ANCount); i++ {
		var a *ResourceRecord
		err = a.unmarshal(r)
		if err != nil {
			return err
		}
		answers = append(answers, *a)
	}

	authorities := []ResourceRecord{}
	for i := 0; i < int(h.NSCount); i++ {
		var a *ResourceRecord
		err = a.unmarshal(r)
		if err != nil {
			return err
		}
		authorities = append(authorities, *a)
	}

	additional := []ResourceRecord{}
	for i := 0; i < int(h.ARCount); i++ {
		var a *ResourceRecord
		err = a.unmarshal(r)
		if err != nil {
			return err
		}
		additional = append(additional, *a)
	}

	p = &Packet{
		Header:      *h,
		Questions:   questions,
		Answers:     answers,
		Authorities: authorities,
		Additional:  additional,
	}

	return nil
}

func (p *Packet) Marshal() ([]byte, error) {
	w := newPacketWriter()

	err := p.Header.marshal(w)
	if err != nil {
		return []byte{}, err
	}

	for i := range p.Questions {
		err = p.Questions[i].marshal(w)
		if err != nil {
			return []byte{}, err
		}
	}

	records := []ResourceRecord{}
	records = append(records, p.Answers...)
	records = append(records, p.Authorities...)
	records = append(records, p.Additional...)

	for i := range records {
		err = records[i].marshal(w)
		if err != nil {
			return []byte{}, err
		}
	}
	return w.buf.Bytes(), nil
}
