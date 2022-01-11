package dns

import (
	"bytes"
	"encoding/binary"
)

type offsetter interface {
	Offset() int
}

type dataWriter interface {
	WriteData(interface{}) error
}

type nameWriter interface {
	WriteName(Name) error
}

type writer interface {
	dataWriter
	nameWriter
}

type offsetLookup struct {
	cache map[string]int
}

func newOffsetLookup() *offsetLookup {
	return &offsetLookup{make(map[string]int)}
}

func (o *offsetLookup) Get(n Name) (int, bool) {
	offset, exists := o.cache[n.LowerString()]
	return offset, exists
}

func (o *offsetLookup) Set(n Name, offset int) {
	o.cache[n.LowerString()] = offset
}

type packetWriter struct {
	buf       *bytes.Buffer
	offsetLkp *offsetLookup
}

func newPacketWriter() *packetWriter {
	return &packetWriter{
		buf:       &bytes.Buffer{},
		offsetLkp: newOffsetLookup(),
	}
}

func (p *packetWriter) WriteData(data interface{}) error {
	return binary.Write(p.buf, binary.BigEndian, data)
}

func (p *packetWriter) WriteName(n Name) error {
	name := make(Name, len(n))
	copy(name, n)
	for {
		if !name.HasParent() {
			break
		}

		initialOffset := p.buf.Len()

		pointerOffset, exists := p.offsetLkp.Get(name)
		if exists {
			pointer := createPointer(pointerOffset)
			return p.WriteData(pointer)
		}

		label := name.LabelAt(0)

		err := p.WriteData(byte(len(label)))
		if err != nil {
			return err
		}

		err = p.WriteData(label)
		if err != nil {
			return err
		}

		p.offsetLkp.Set(name, initialOffset)

		name = name.Parent()
	}

	// If there was a pointer it would return early
	// so this will only occcur if it is not capable
	// of being compressed
	return p.WriteData(nameTerminator)
}
