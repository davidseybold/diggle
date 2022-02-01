package dns

import (
	"encoding/binary"
)

const (
	pointerMask byte = 0xC0
	labelMask   byte = 0x3F
)

func createPointer(offset int) uint16 {
	msb, lsb := byte(offset>>8), byte(offset)
	return binary.BigEndian.Uint16([]byte{(pointerMask | msb), lsb})
}

func createPointerOffset(firstOctet byte, secondOctet byte) uint16 {
	numFirst := labelMask & firstOctet
	return binary.BigEndian.Uint16([]byte{numFirst, secondOctet})
}

func isPointerSignal(p byte) bool {
	return p >= 192
}
