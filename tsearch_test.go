package tsearch

import (
	"fmt"
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

func TestSimilarity(t *testing.T) {
	a := "niu lan qi an quan"
	b := "wei ruan bi ying liu lan qi an quan fang hu"
	separator := NewTrigramSeparator()
	textSearch := NewTextSearch(separator, nil)

	value := textSearch.WordSimilarity(a, b)
	fmt.Println(value)

}
