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

func (q *Question) encode(w writeOffsetter, c *compressionCache) error {
	if err := q.Name.encode(w, c); err != nil {
		return err
	}

	if err := q.Type.encode(w, c); err != nil {
		return err
	}

	return q.Class.encode(w, c)
}

func (q *Question) decode(r readSeekOffsetter) error {

	if err := q.Name.decode(r); err != nil {
		return err
	}

	if err := q.Type.decode(r); err != nil {
		return err
	}

	return q.Class.decode(r)
}
