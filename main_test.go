package main

import (
	"bytes"
	"testing"
)

func TestWriteAliases(t *testing.T) {
	cases := []struct {
		aliases aliases
		want    string
	}{
		{
			map[string][]string{
				"testing": []string{
					"walter",
					"mike",
					"donald",
				},
				"also_testing": []string{
					"joe",
					"woody",
					"theo",
					"jack",
				},
			},
			"also_testing: joe, woody, theo, jack\ntesting: walter, mike, donald\n",
		},
		{
			map[string][]string{
				"testing": []string{
					"walter",
					"bob",
					"donald",
				},
			},
			"testing: walter, bob, donald\n",
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
