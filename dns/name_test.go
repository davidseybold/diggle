package dns

import (
	"testing"
)

func TestNameEquals(t *testing.T) {
	cases := []struct {
		name     string
		a        Name
		b        Name
		expected bool
	}{
		{
			name:     "all same case are equal",
			a:        Name{{'g', 'o', 'o', 'g', 'l', 'e'}, {'c', 'o', 'm'}},
			b:        Name{{'g', 'o', 'o', 'g', 'l', 'e'}, {'c', 'o', 'm'}},
			expected: true,
		},
		{
			name:     "different cases are equal",
			a:        Name{{'g', 'o', 'o', 'g', 'L', 'e'}, {'c', 'o', 'M'}},
			b:        Name{{'g', 'o', 'O', 'g', 'l', 'e'}, {'c', 'o', 'm'}},
			expected: true,
		},
		{
			name:     "different domains are not equal",
			a:        Name{{'c', 'o', 'M'}},
			b:        Name{{'g', 'o', 'O', 'g', 'l', 'e'}, {'c', 'o', 'm'}},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.a.Equals(c.b) != c.expected {
				t.FailNow()
			}
		})
	}
}

func TestNameParent(t *testing.T) {
	cases := []struct {
		name     string
		domain   Name
		expected Name
	}{
		{
			name:     "correct parent is returned for two level",
			domain:   Name{{'g', 'o', 'o', 'g', 'l', 'e'}, {'c', 'o', 'm'}},
			expected: Name{{'c', 'o', 'm'}},
		},
		{
			name:     "correct parent is returned for one level",
			domain:   Name{{'c', 'o', 'm'}},
			expected: Name{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			parent := c.domain.Parent()
			if !parent.Equals(c.expected) {
				t.Fatal("expected", c.expected, "but got", parent)
			}
		})
	}
}

func TestNameHasParent(t *testing.T) {
	cases := []struct {
		name     string
		domain   Name
		expected bool
	}{
		{
			name:     "two level domain has a parent",
			domain:   Name{{'g', 'o', 'o', 'g', 'l', 'e'}, {'c', 'o', 'm'}},
			expected: true,
		},
		{
			name:     "one level domain has a parent",
			domain:   Name{{'c', 'o', 'm'}},
			expected: true,
		},
		{
			name:     "root domain has no parent",
			domain:   Name{},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.domain.HasParent() != c.expected {
				t.FailNow()
			}
		})
	}
}
