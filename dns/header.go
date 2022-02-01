package dns

const (
	maskQR     uint16 = 1 << 15
	maskOpcode uint16 = 0x7800
	maskAA     uint16 = 1 << 10
	maskTC     uint16 = 1 << 9
	maskRD     uint16 = 1 << 8
	maskRA     uint16 = 1 << 7
	maskAD     uint16 = 1 << 5
	maskCD     uint16 = 1 << 4
	maskZ      uint16 = 0x70
	maskRCode  uint16 = 0xF

	opcodeShift = 11

	headerLen = 12
)

type Flags struct {
	AuthoritativeAnswer bool
	Truncated           bool
	RecursionDesired    bool
	RecursionAvailable  bool
	AuthenticData       bool
	CheckingDisabled    bool
}

type header struct {
	// A 16 bit identifier assigned by the program that
	// generates any kind of query.  This identifier is copied
	// the corresponding reply and can be used by the requester
	// to match up replies to outstanding queries.
	ID uint16
	// A one bit field that specifies whether this message is a
	// query (0), or a response (1)
	Type   bool
	Opcode byte
	Flags  Flags

	ResponseCode ResponseCode

	questionCount   uint16
	answerCount     uint16
	authorityCount  uint16
	additionalCount uint16
}

func (h *header) encode(w writeOffsetter, c *compressionCache) error {
	if err := writeUint16(w, h.ID); err != nil {
		return err
	}

	b := addBits(0, encodeBool(h.Type), 1)
	b = addBits(b, h.Opcode, 4)
	b = addBits(b, encodeBool(h.Flags.AuthoritativeAnswer), 1)
	b = addBits(b, encodeBool(h.Flags.Truncated), 1)
	b = addBits(b, encodeBool(h.Flags.RecursionDesired), 1)

	if err := writeByte(w, b); err != nil {
		return err
	}

	b = addBits(0, encodeBool(h.Flags.RecursionAvailable), 1)
	b = addBits(b, 0, 1) // Write reserved bit set to 0
	b = addBits(b, encodeBool(h.Flags.AuthenticData), 1)
	b = addBits(b, encodeBool(h.Flags.CheckingDisabled), 1)
	b = addBits(b, byte(h.ResponseCode), 4)

	if err := writeByte(w, b); err != nil {
		return err
	}

	counts := []uint16{h.questionCount, h.answerCount, h.authorityCount, h.additionalCount}

	for i := range counts {
		if err := writeUint16(w, counts[i]); err != nil {
			return err
		}
	}
	return nil
}

func addBits(dst byte, src byte, numBits int) byte {
	return (dst << numBits) & src
}

func encodeBool(b bool) byte {
	var val byte
	if b {
		val = 1
	}
	return val
}

func (h *header) decode(r readSeekOffsetter) error {
	headerSections := make([]uint16, 6)
	for i := range headerSections {
		section, err := readUint16(r)
		if err != nil {
			return err
		}
		headerSections[i] = section
	}
	metadata := headerSections[1]

	h = &header{
		ID:     headerSections[0],
		Type:   getQR(metadata),
		Opcode: getOpCode(metadata),
		Flags: Flags{
			AuthoritativeAnswer: getAA(metadata),
			Truncated:           getTC(metadata),
			RecursionDesired:    getRD(metadata),
			RecursionAvailable:  getRA(metadata),
			AuthenticData:       getAD(metadata),
			CheckingDisabled:    getCD(metadata),
		},
		ResponseCode:    ResponseCode(getRCode(metadata)),
		questionCount:   headerSections[2],
		answerCount:     headerSections[3],
		authorityCount:  headerSections[4],
		additionalCount: headerSections[5],
	}
	return nil
}

func getQR(d uint16) bool {
	return (d & maskQR) > 0
}

func getOpCode(d uint16) byte {
	return byte((d & maskOpcode) >> opcodeShift)
}

func getAA(d uint16) bool {
	return (d & maskAA) > 0
}

func getTC(d uint16) bool {
	return (d & maskTC) > 0
}

func getRD(d uint16) bool {
	return (d & maskRD) > 0
}

func getRA(d uint16) bool {
	return (d & maskRA) > 0
}

func getAD(d uint16) bool {
	return (d & maskAD) > 0
}

func getCD(d uint16) bool {
	return (d & maskCD) > 0
}

func getRCode(flags uint16) uint8 {
	return uint8(flags & maskRCode)
}
