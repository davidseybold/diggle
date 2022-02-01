package dns

type unsupportedRecordData struct {
	rdLength uint16
	rData    []byte
}

func (u *unsupportedRecordData) encode(w writeOffsetter, c *compressionCache) error {
	if err := writeUint16(w, u.rdLength); err != nil {
		return err
	}

	_, err := w.Write(u.rData)
	return err
}

func (u *unsupportedRecordData) decode(r readSeekOffsetter) error {
	var err error
	if u.rdLength, err = readUint16(r); err != nil {
		return err
	}

	if u.rData, err = readNBytes(r, int(u.rdLength)); err != nil {
		return err
	}

	return nil
}
