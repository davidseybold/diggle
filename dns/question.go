package dns

type Question struct {
	// a domain name represented as a sequence of labels, where
	// each label consists of a length octet followed by that
	// number of octets.  The domain name terminates with the
	// zero length octet for the null label of the root.  Note
	// that this field may be an odd number of octets; no
	// padding is used.
	Name Name
	// a two octet code which specifies the type of the query.
	// The values for this field include all codes valid for a
	// TYPE field, together with some more general codes which
	// can match more than one type of RR.
	Type Type
	// a two octet code that specifies the class of the query.
	// For example, the QCLASS field is IN for the Internet.
	Class Class
}

func (q *Question) marshal(w writer) error {
	err := w.WriteName(q.Name)
	if err != nil {
		return err
	}

	err = w.WriteData(uint16(q.Type))
	if err != nil {
		return err
	}

	return w.WriteData(uint16(q.Class))
}

func (q *Question) unmarshal(r reader) error {
	name, err := r.ReadName()
	if err != nil {
		return err
	}

	qType, err := r.ReadUint16()
	if err != nil {
		return err
	}

	qClass, err := r.ReadUint16()
	if err != nil {
		return err
	}

	q = &Question{
		Name:  name,
		Type:  Type(qType),
		Class: Class(qClass),
	}
	return nil
}
