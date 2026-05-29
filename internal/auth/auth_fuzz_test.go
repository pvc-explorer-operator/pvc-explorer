package auth

import (
	"strings"
	"testing"
)

func FuzzSplitTrim(f *testing.F) {
	seeds := []string{
		"",
		",",
		"admin",
		"admin,viewer",
		" admin , viewer ",
		"a,,b, , c",
		"\nadmin,\tviewer",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		out := splitTrim(s)
		for _, v := range out {
			if v == "" {
				t.Fatalf("splitTrim returned empty token for input %q", s)
			}
			if strings.TrimSpace(v) != v {
				t.Fatalf("splitTrim returned untrimmed token %q for input %q", v, s)
			}
		}
	})
}
