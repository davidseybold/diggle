package dns

import (
	"bytes"
	"io"
)

type decoder interface {
	decode(readSeekOffsetter) error
}

type readSeekOffsetter interface {
	io.ReadSeeker
	Offset() int
}

type offsetReader struct {
	*bytes.Reader
}

func newOffsetReader(buf []byte) *offsetReader {
	return &offsetReader{
		Reader: bytes.NewReader(buf),
	}
}

func (p *offsetReader) Offset() int {
	return int(p.Size()) - p.Len()
}

func DecodePacket(b []byte) (Packet, error) {
	r := newOffsetReader(b)

	var h *header
	if err := h.decode(r); err != nil {
		return Packet{}, err
	}

	questions := []Question{}
	for i := 0; i < int(h.questionCount); i++ {
		var q *Question
		if err := q.decode(r); err != nil {
			return Packet{}, err
		}
		questions = append(questions, *q)
	}

	answers := make([]ResourceRecord, h.answerCount)
	authorities := make([]ResourceRecord, h.authorityCount)
	additional := make([]ResourceRecord, h.additionalCount)

	tRecs := h.answerCount + h.authorityCount + h.additionalCount

	recs := []ResourceRecord{}
	for i := 0; i < int(tRecs); i++ {
		var rr *ResourceRecord
		if err := rr.decode(r); err != nil {
			return Packet{}, err
		}
		recs = append(recs, *rr)
	}

	copy(answers, recs[:h.answerCount])
	copy(authorities, recs[h.answerCount:h.answerCount+h.authorityCount])
	copy(additional, recs[h.answerCount+h.authorityCount:])

	return Packet{
		header:      *h,
		Questions:   questions,
		Answers:     answers,
		Authorities: authorities,
		Additional:  additional,
	}, nil
}
