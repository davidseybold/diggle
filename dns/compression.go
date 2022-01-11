package dns

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	pointerMask byte = 0xC0
	labelMask   byte = 0x3F
)

func decompress(r io.ByteReader, firstOctet byte, domainLkp nameLookup) (Name, error) {
	secondOctet, err := r.ReadByte()
	if err != nil {
		return Name{}, err
	}
	pointerOffset := createPointerOffset(firstOctet, secondOctet)

	pointer, exists := domainLkp[pointerOffset]
	if !exists {
		return Name{}, errors.New("invalid pointer")
	}

	return pointer, nil
}

func createPointer(pointerOffset int) uint16 {
	msb, lsb := byte(pointerOffset>>8), byte(pointerOffset)
	return binary.BigEndian.Uint16([]byte{(pointerMask | msb), lsb})
}

func createPointerOffset(firstOctet byte, secondOctet byte) int {
	numFirst := labelMask & firstOctet
	return int(binary.BigEndian.Uint16([]byte{numFirst, secondOctet}))
}

func isPointerSignal(p byte) bool {
	return p >= 192
}
