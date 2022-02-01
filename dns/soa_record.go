package dns

import "fmt"

type SOARecordData struct {
	MName   Name
	RName   Name
	Serial  uint32
	Refresh int32
	Retry   int32
	Expire  int32
	Minimum uint32
}

func (s SOARecordData) String() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", s.MName, s.RName, s.Serial, s.Refresh, s.Retry, s.Expire, s.Minimum)
}

func (s SOARecordData) encode(w writeOffsetter, c *compressionCache) error {
	buf := newOffsetWriter(w.Offset())
	if err := s.MName.encode(buf, c); err != nil {
		return err
	}

	if err := s.RName.encode(buf, c); err != nil {
		return err
	}

	if err := writeUint32(buf, s.Serial); err != nil {
		return err
	}

	if err := writeInt32(buf, s.Refresh); err != nil {
		return err
	}

	if err := writeInt32(buf, s.Retry); err != nil {
		return err
	}

	if err := writeInt32(buf, s.Expire); err != nil {
		return err
	}

	if err := writeUint32(buf, s.Minimum); err != nil {
		return err
	}

	if err := writeUint32(w, buf.Len()); err != nil {
		return err
	}

	_, err := buf.WriteTo(buf)
	if err != nil {
		return err
	}

	return nil
}

func (s *SOARecordData) decode(r readSeekOffsetter) error {

	// Discard RDLength as we know all the data
	_, err := readUint32(r)
	if err != nil {
		return err
	}

	if err := s.MName.decode(r); err != nil {
		return err
	}

	if err := s.RName.decode(r); err != nil {
		return err
	}

	if s.Serial, err = readUint32(r); err != nil {
		return err
	}

	if s.Refresh, err = readInt32(r); err != nil {
		return err
	}

	if s.Retry, err = readInt32(r); err != nil {
		return err
	}

	if s.Expire, err = readInt32(r); err != nil {
		return err
	}

	if s.Minimum, err = readUint32(r); err != nil {
		return err
	}

	return nil
}
