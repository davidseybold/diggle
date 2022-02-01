package dns

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	nameTerminator byte = 0
)

type Name [][]byte

func NewName(s string) Name {
	if s[len(s)-1] != '.' {
		s = s + "."
	}
	return bytes.Split([]byte(s), []byte{'.'})
}

func (n Name) Parent() Name {
	if len(n) == 0 {
		return Name{{}}
	}
	p := n[1:]
	pCopy := make(Name, len(p))
	copy(pCopy, p)
	return pCopy
}

func (n Name) LabelAt(i int) []byte {
	return n[i]
}

func (n Name) HasParent() bool {
	return len(n) > 0
}

func (n Name) Equals(x Name) bool {
	return n.LowerString() == x.LowerString()
}

func (n Name) String() string {
	return string(bytes.Join(n, []byte{'.'}))
}

func (n Name) LowerString() string {
	return strings.ToLower(n.String())
}

func (n Name) encode(w writeOffsetter, c *compressionCache) error {
	name := make(Name, len(n))
	copy(name, n)
	for {
		if !name.HasParent() {
			break
		}

		initialOffset := w.Offset()

		pointerOffset, exists := c.Get(name)
		if exists {
			pointer := createPointer(pointerOffset)
			return writeUint16(w, pointer)
		}

		label := name.LabelAt(0)

		if err := writeByte(w, len(label)); err != nil {
			return err
		}

		_, err := w.Write(label)
		if err != nil {
			return err
		}

		c.Set(name, initialOffset)

		name = name.Parent()
	}
	return nil
}

func (n *Name) decode(r readSeekOffsetter) error {
	var decodeLabels func(Name, int) (Name, error)
	decodeLabels = func(name Name, ptrCnt int) (Name, error) {

		b, err := readByte(r)
		if err != nil {
			return Name{}, err
		}

		if isNameTerminator(b) {
			name = append(name, []byte{0})
			return name, nil
		} else if isLabelSignal(b) {
			label, err := readNBytes(r, int(b))
			if err != nil {
				return Name{}, err
			}
			name = append(name, label)
			return decodeLabels(name, ptrCnt)
		} else if isPointerSignal(b) {
			if ptrCnt == 0 {
				return Name{}, errors.New("invalid packet: multiple sequential pointers")
			}
			secondOctet, err := readByte(r)
			if err != nil {
				return Name{}, err
			}
			ptrOffset := createPointerOffset(b, secondOctet)
			crnt := r.Offset()

			_, err = r.Seek(int64(ptrOffset), io.SeekStart)
			if err != nil {
				return Name{}, err
			}

			name, err := decodeLabels(name, ptrCnt-1)
			if err != nil {
				return Name{}, err
			}

			_, err = r.Seek(int64(crnt), io.SeekStart)
			if err != nil {
				return Name{}, err
			}

			return name, nil

		} else {
			return Name{}, errors.New("invalid byte encountered")
		}
	}

	name, err := decodeLabels(Name{}, 1)

	*n = name

	return err
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
