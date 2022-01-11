package dns

import (
	"bytes"
	"strings"
)

const (
	nameTerminator byte = 0x0
)

type Name [][]byte

func (d Name) Parent() Name {
	if len(d) == 0 {
		return Name{{}}
	}
	p := d[1:]
	pCopy := make(Name, len(p))
	copy(pCopy, p)
	return pCopy
}

func (d Name) LabelAt(i int) []byte {
	return d[i]
}

func (d Name) HasParent() bool {
	return len(d) > 0
}

func (d Name) Equals(x Name) bool {
	return d.LowerString() == x.LowerString()
}

func (d Name) String() string {
	return string(d.join())
}

func (d Name) LowerString() string {
	return strings.ToLower(d.String())
}

func (d Name) join() []byte {
	return append(bytes.Join(d, []byte{'.'}), '.')
}

func isNameTerminator(b byte) bool {
	return b == nameTerminator
}

func isLabelSignal(p byte) bool {
	return p <= 63
}
