package dns

import (
	"errors"
	"net"
)

// An A record
type ARecordData struct {
	Address net.IP
}

func (a ARecordData) encode(w writeOffsetter, c *compressionCache) error {
	err := writeByte(w, byte(4))
	if err != nil {
		return err
	}
	_, err = w.Write(a.Address.To4())
	return err
}

func (a *ARecordData) decode(r offsetReader) error {
	dLen, err := readUint16(r)
	if err != nil {
		return err
	}

	if dLen != 4 {
		return errors.New("invalid a record")
	}

	d, err := readNBytes(r, int(dLen))
	if err != nil {
		return err
	}

	a.Address = net.IP(d)

	return nil
}

func (a ARecordData) String() string {
	return a.Address.String()
}
