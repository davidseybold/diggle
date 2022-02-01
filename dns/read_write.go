package dns

import (
	"encoding/binary"
	"errors"
	"io"
)

func writeData(w io.Writer, d interface{}) error {
	return binary.Write(w, binary.BigEndian, d)
}

func writeByte(w io.Writer, d interface{}) error {
	n, ok := d.(byte)
	if !ok {
		return errors.New("cannot convert to byte")
	}
	return writeData(w, n)
}

func writeUint8(w io.Writer, d interface{}) error {
	n, ok := d.(uint8)
	if !ok {
		return errors.New("cannot convert to uint8")
	}
	return writeData(w, n)
}

func writeUint16(w io.Writer, d interface{}) error {
	n, ok := d.(uint16)
	if !ok {
		return errors.New("cannot convert to uint16")
	}
	return writeData(w, n)
}

func writeInt32(w io.Writer, d interface{}) error {
	n, ok := d.(int32)
	if !ok {
		return errors.New("cannot convert to int32")
	}
	return writeData(w, n)
}

func writeUint32(w io.Writer, d interface{}) error {
	n, ok := d.(uint32)
	if !ok {
		return errors.New("cannot convert to uint32")
	}
	return writeData(w, n)
}

func writeUint64(w io.Writer, d interface{}) error {
	n, ok := d.(uint64)
	if !ok {
		return errors.New("cannot convert to uint64")
	}
	return writeData(w, n)
}

func readData(r io.Reader, d interface{}) error {
	return binary.Read(r, binary.BigEndian, d)
}

func readByte(r io.Reader) (byte, error) {
	var b byte
	err := readData(r, &b)
	return b, err
}

func readNBytes(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	err := readData(r, b)
	return b, err
}

func readUint8(r io.Reader) (uint8, error) {
	var b uint8
	err := readData(r, &b)
	return b, err
}

func readUint16(r io.Reader) (uint16, error) {
	var b uint16
	err := readData(r, &b)
	return b, err
}

func readUint32(r io.Reader) (uint32, error) {
	var b uint32
	err := readData(r, &b)
	return b, err
}

func readInt32(r io.Reader) (int32, error) {
	var b int32
	err := readData(r, &b)
	return b, err
}

func readUint64(r io.Reader) (uint64, error) {
	var b uint64
	err := readData(r, &b)
	return b, err
}
