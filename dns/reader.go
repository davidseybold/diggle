package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type dataReader interface {
	ReadData(interface{}) error
}

type nameReader interface {
	ReadName() (Name, error)
}

type uintReader interface {
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
}

type byteReader interface {
	io.ByteReader
	ReadNBytes(int) ([]byte, error)
}

type reader interface {
	byteReader
	dataReader
	nameReader
	uintReader
}

type nameLookup map[int]Name

type packetReader struct {
	*bytes.Reader
	nameLkp nameLookup
}

func newPacketReader(buf []byte) *packetReader {
	return &packetReader{
		Reader:  bytes.NewReader(buf),
		nameLkp: make(nameLookup),
	}
}

func (p *packetReader) Offset() int {
	return int(p.Size()) - p.Len()
}

func (b *packetReader) ReadUint16() (uint16, error) {
	var n uint16
	err := b.ReadData(&n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (p *packetReader) ReadUint32() (uint32, error) {
	var n uint32
	err := p.ReadData(&n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (p *packetReader) ReadNBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	err := p.ReadData(buf)
	if err != nil {
		return []byte{}, err
	}

	return buf, nil
}

func (p *packetReader) ReadData(data interface{}) error {
	return binary.Read(p, binary.BigEndian, data)
}

func (p *packetReader) ReadName() (Name, error) {
	beginningOffset := p.Offset()

	var decodeLabels func(Name) (Name, error)
	decodeLabels = func(name Name) (Name, error) {

		b, err := p.ReadByte()
		if err != nil {
			return Name{}, err
		}

		switch {
		case isNameTerminator(b):
			return name, nil
		case isLabelSignal(b):
			label, err := p.ReadNBytes(int(b))
			if err != nil {
				return Name{}, err
			}
			name = append(name, label)
			name, err = decodeLabels(name)
			if err != nil {
				return Name{}, err
			}
			return name, nil
		case isPointerSignal(b):
			pointer, err := decompress(p, b, p.nameLkp)
			if err != nil {
				return Name{}, err
			}
			return append(name, pointer...), nil
		default:
			return Name{}, errors.New("invalid byte encountered")
		}
	}

	name, err := decodeLabels(Name{})
	if err != nil {
		return Name{}, err
	}

	p.nameLkp[beginningOffset] = name

	return Name{}, nil
}
