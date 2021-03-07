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
		str    string
		tokens []string
	}

	tests := []Test{
		{
			"word",
			[]string{"  w", " wo", "wor", "ord", "rd "},
		},
		{
			"two words",
			[]string{"  t", " tw", "two", "wo ", "  w", " wo", "wor", "ord", "rds", "ds "},
		},
	}

	for i, test := range tests {
		output := separator.Extract(test.str)
		if !cmpStringArray(test.tokens, output) {
			t.Errorf("Extract error %d\n", i)
		}
	}
}

func TestDistinctTokens(t *testing.T) {
	type Test struct {
		tokens []string
		result []string
	}

	tests := []Test{
		{
			[]string{" w", " wo", "wor", "ord", "rd "},
			[]string{" w", " wo", "wor", "ord", "rd "},
		},
		{
			[]string{"a", "b", "a", "b"},
			[]string{"a", "b"},
		},
	}

	for i, test := range tests {
		output := DistinctTokens(test.tokens)
		if !cmpStringArray(test.result, output) {
			t.Errorf("DistinctTokens error %d\n", i)
		}
	}
}
