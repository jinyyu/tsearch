package tsearch

type Separator interface {
	Extract(str string) []string
}

type trigramSeparator struct {
	numberCharacters int
}

func NewTrigramSeparator() Separator {
	return &trigramSeparator{
		numberCharacters: 3,
	}
}

func (t *trigramSeparator) Extract(str string) []string {
	str = "  " + str + " "
	ret := make([]string, len(str)-t.numberCharacters+1)
	for i := 0; i < len(str)-t.numberCharacters+1; i++ {
		word := str[i : i+t.numberCharacters]
		ret[i] = word
	}
	return ret
}
