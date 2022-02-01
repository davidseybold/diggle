package dns

import (
	"bytes"
	"io"
)

type shouldWriteFunc func(*offsetWriter, *offsetWriter) bool

func shouldWriteUDP(w, body *offsetWriter) bool {
	return w.Len()+body.Len() > maxUDPPacketSize
}

func alwaysWrite(_, _ *offsetWriter) bool {
	return true
}

func EncodeTCPPacket(p Packet) (EncodePacketResult, error) {
	return encode(p, alwaysWrite)
}

func EncodeUDPPacket(p Packet) (EncodePacketResult, error) {
	return encode(p, shouldWriteUDP)
}

func encode(p Packet, shouldWrite shouldWriteFunc) (EncodePacketResult, error) {
	body := newOffsetWriter(headerLen)
	cache := newCompressionCache()
	trunc := false

	var qCnt int
	for _, q := range p.Questions {
		qw := newOffsetWriter(body.Offset())
		err := q.encode(qw, cache)
		if err != nil {
			return EncodePacketResult{}, err
		}
		if shouldWriteUDP(qw, body) {
			cache.Revert()
			trunc = true
			continue
		}
		cache.Commit()
		_, err = qw.WriteTo(body)
		if err != nil {
			return EncodePacketResult{}, err
		}
		qCnt++
	}

	anCnt, err := encodeRR(body, cache, shouldWrite, p.Answers)
	if err != nil {
		return EncodePacketResult{}, err
	}
	if anCnt < len(p.Answers) {
		trunc = true
	}

	nsCnt, err := encodeRR(body, cache, shouldWrite, p.Authorities)
	if err != nil {
		return EncodePacketResult{}, err
	}
	if nsCnt < len(p.Authorities) {
		trunc = true
	}

	adCnt, err := encodeRR(body, cache, shouldWrite, p.Additional)
	if err != nil {
		return EncodePacketResult{}, err
	}
	if adCnt < len(p.Additional) {
		trunc = true
	}

	p.header.questionCount = uint16(qCnt)
	p.header.answerCount = uint16(anCnt)
	p.header.authorityCount = uint16(nsCnt)
	p.header.additionalCount = uint16(adCnt)
	p.header.Flags.Truncated = trunc

	buf := newOffsetWriter(0)
	if err := p.header.encode(buf, cache); err != nil {
		return EncodePacketResult{}, err
	}

	_, err = body.WriteTo(buf)
	if err != nil {
		return EncodePacketResult{}, err
	}

	return EncodePacketResult{
		Bytes:     buf.Bytes(),
		Truncated: trunc,
	}, nil
}

func encodeRR(body *offsetWriter, c *compressionCache, shouldWrite shouldWriteFunc, rr []ResourceRecord) (int, error) {
	var cnt int
	for _, r := range rr {
		w := newOffsetWriter(body.Offset())
		err := r.encode(w, c)
		if err != nil {
			return cnt, err
		}
		if shouldWrite(w, body) {
			c.Revert()
			continue
		}
		c.Commit()
		_, err = w.WriteTo(body)
		if err != nil {
			return cnt, err
		}
		cnt++
	}
	return cnt, nil
}

type writeOffsetter interface {
	io.Writer
	Offset() int
}

type encoder interface {
	encode(writeOffsetter, *compressionCache) error
}

type EncodePacketResult struct {
	Bytes     []byte
	Truncated bool
}

type compressionCache struct {
	cache     map[string]int
	lastAdded []string
}

func newCompressionCache() *compressionCache {
	return &compressionCache{
		cache:     make(map[string]int),
		lastAdded: []string{},
	}
}

func (c *compressionCache) Set(n Name, i int) {
	c.cache[n.LowerString()] = i
	c.lastAdded = append(c.lastAdded, n.LowerString())
}

func (c *compressionCache) Get(n Name) (int, bool) {
	o, exists := c.cache[n.LowerString()]
	return o, exists
}

func (c *compressionCache) Revert() {
	for _, k := range c.lastAdded {
		delete(c.cache, k)
	}
}

func (c *compressionCache) Commit() {
	c.lastAdded = []string{}
}

type offsetWriter struct {
	bytes.Buffer
	offset int
}

func newOffsetWriter(i int) *offsetWriter {
	return &offsetWriter{
		offset: i,
	}
}

func (o *offsetWriter) Write(b []byte) (int, error) {
	n, err := o.Buffer.Write(b)
	if err == nil {
		o.offset += n
	}
	return n, err
}

func (o *offsetWriter) Offset() int {
	return o.offset
}
