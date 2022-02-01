package dns

import (
	"errors"
	"net"
)

type AAAARecordData struct {
	Address net.IP
}

func (a AAAARecordData) encode(w writeOffsetter, c *compressionCache) error {
	err := writeByte(w, byte(16))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(a.Address))
	return err
}

func (a *AAAARecordData) decode(r offsetReader) error {
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

func (a AAAARecordData) String() string {
	return a.Address.String()
}
