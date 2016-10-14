package main

import (
	"bytes"
	"testing"
)

func TestWriteAliases(t *testing.T) {
	cases := []struct {
		aliases []alias
		want    string
	}{
		{
			[]alias{
				alias{
					"testing",
					[]string{
						"donald",
						"walter",
						"bob",
					},
				},
			},
			"testing: donald, walter, bob\n",
		},
	}
	var b bytes.Buffer

	for _, c := range cases {
		b.Reset()
		err := writeAliases(c.aliases, &b)
		if err != nil {
			t.Error(err)
		}
		got := b.String()
		if got != c.want {
			t.Errorf("writeAliases failure: want %s, got %s", c.want, got)
		}
	}
}
