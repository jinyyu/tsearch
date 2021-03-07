package tsearch

import (
	"strings"
	"testing"
)

func cmpStringArray(a []string, b []string) bool {
	aStr := strings.Join(a, " ")
	bStr := strings.Join(b, " ")
	return aStr == bStr
}

func Test_TrigramSeparator(t *testing.T) {
	separator := NewTrigramSeparator()

	type Test struct {
		str   string
		array []string
	}

	tests := []Test{
		{
			"word",
			[]string{" w", " wo", "wor", "ord", "rd "},
		},
		{
			"two words",
			[]string{" t", " tw", "two", "wo ", " w", " wo", "wor", "ord", "rds", "ds "},
		},
	}

	for i, test := range tests {
		output := separator.Extract(test.str)
		if cmpStringArray(test.array, output) {
			t.Errorf("Extract error %d\n", i)
		}
	}
}
