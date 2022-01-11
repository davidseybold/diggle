package dns

const (
	maskQR     uint16 = 0x8000
	maskOpcode uint16 = 0x7800
	maskAA     uint16 = 0x400
	maskTC     uint16 = 0x200
	maskRD     uint16 = 0x100
	maskRA     uint16 = 0x80
	maskZ      uint16 = 0x70
	maskRCode  uint16 = 0xF

	shiftQR     = 15
	shiftOpcode = 11
	shiftAA     = 10
	shiftTC     = 9
	shiftRD     = 8
	shiftRA     = 7
	shiftZ      = 4
)

type Header struct {
	// A 16 bit identifier assigned by the program that
	// generates any kind of query.  This identifier is copied
	// the corresponding reply and can be used by the requester
	// to match up replies to outstanding queries.
	ID uint16
	// A one bit field that specifies whether this message is a
	// query (0), or a response (1)
	QR bool
	// A four bit field that specifies kind of query in this
	// message.  This value is set by the originator of a query
	// and copied into the response.  The values are:
	// 0               a standard query (QUERY)
	// 1               an inverse query (IQUERY)
	// 2               a server status request (STATUS)
	// 3-15            reserved for future use
	Opcode byte
	// Authoritative Answer - this bit is valid in responses,
	// and specifies that the responding name server is an
	// authority for the domain name in question section.

	// Note that the contents of the answer section may have
	// multiple owner names because of aliases.  The AA bit
	// corresponds to the name which matches the query name, or
	// the first owner name in the answer section.
	AA bool
	// TrunCation - specifies that this message was truncated
	// due to length greater than that permitted on the
	// transmission channel
	TC bool
	// Recursion Desired - this bit may be set in a query and
	// is copied into the response.  If RD is set, it directs
	// the name server to pursue the query recursively.
	// Recursive query support is optional.
	RD bool
	// Recursion Available - this be is set or cleared in a
	// response, and denotes whether recursive query support is
	// available in the name server.
	RA bool
	// 3 bits Reserved for future use.  Must be zero in all queries
	// and responses.
	Z byte
	// Response code - this 4 bit field is set as part of
	// responses.  The values have the following
	// interpretation:

	// 0               No error condition

	// 1               Format error - The name server was
	// 								unable to interpret the query.

	// 2               Server failure - The name server was
	// 								unable to process this query due to a
	// 								problem with the name server.

	// 3               Name Error - Meaningful only for
	// 								responses from an authoritative name
	// 								server, this code signifies that the
	// 								domain name referenced in the query does
	// 								not exist.

	// 4               Not Implemented - The name server does
	// 								not support the requested kind of query.

	// 5               Refused - The name server refuses to
	// 								perform the specified operation for
	// 								policy reasons.  For example, a name
	// 								server may not wish to provide the
	// 								information to the particular requester,
	// 								or a name server may not wish to perform
	// 								a particular operation (e.g., zone
	// transfer) for particular data.
	RCode uint8
	// an unsigned 16 bit integer specifying the number of
	// entries in the question section.
	QDCount uint16
	// an unsigned 16 bit integer specifying the number of
	// resource records in the answer section.
	ANCount uint16
	// an unsigned 16 bit integer specifying the number of name
	// server resource records in the authority records
	// section.
	NSCount uint16
	//an unsigned 16 bit integer specifying the number of
	// resource records in the additional records section.
	ARCount uint16
}

func (h *Header) marshal(w writer) error {
	err := w.WriteData(h.ID)
	if err != nil {
		return err
	}

	b := addBits(0, encodeBool(h.QR), 1)
	b = addBits(b, h.Opcode, 4)
	b = addBits(b, encodeBool(h.AA), 1)
	b = addBits(b, encodeBool(h.TC), 1)
	b = addBits(b, encodeBool(h.RD), 1)

	err = w.WriteData(b)
	if err != nil {
		return err
	}

	c := addBits(0, encodeBool(h.RA), 1)
	c = addBits(c, h.Z, 3)
	c = addBits(c, h.RCode, 4)

	err = w.WriteData(c)
	if err != nil {
		return err
	}

	counts := []uint16{h.QDCount, h.ANCount, h.NSCount, h.ARCount}

	for i := range counts {
		err = w.WriteData(counts[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Header) unmarshal(r reader) error {
	headerSections := make([]uint16, 6)
	for i := range headerSections {
		section, err := r.ReadUint16()
		if err != nil {
			return err
		}
		headerSections[i] = section
	}
	flags := headerSections[1]

	h = &Header{
		ID:      headerSections[0],
		QR:      getQR(flags),
		Opcode:  getOpCode(flags),
		AA:      getAA(flags),
		TC:      getTC(flags),
		RD:      getRD(flags),
		RA:      getRA(flags),
		Z:       getZ(flags),
		RCode:   getRCode(flags),
		QDCount: headerSections[2],
		ANCount: headerSections[3],
		NSCount: headerSections[4],
		ARCount: headerSections[5],
	}
	return nil
}

func getQR(flags uint16) bool {
	return (flags&maskQR)>>shiftQR != 0
}

func getOpCode(flags uint16) byte {
	return byte((flags & maskOpcode) >> shiftOpcode)
}

func getAA(flags uint16) bool {
	return (flags&maskAA)>>shiftAA != 0
}

func getTC(flags uint16) bool {
	return (flags&maskTC)>>shiftTC != 0
}

func getRD(flags uint16) bool {
	return (flags&maskRD)>>shiftRD != 0
}

func getRA(flags uint16) bool {
	return (flags&maskRA)>>shiftRA != 0
}

func getZ(flags uint16) byte {
	return byte((flags & maskZ) >> shiftZ)
}

func getRCode(flags uint16) uint8 {
	return uint8(flags & maskRCode)
}

func encodeBool(b bool) byte {
	var val byte
	if b {
		val = 1
	}
	return val
}

func addBits(dst byte, src byte, numBits int) byte {
	return (dst << numBits) & src
}
