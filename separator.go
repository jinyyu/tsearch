package tsearch

import "strings"

type Separator interface {
	Extract(text string) []string
}

type trigramSeparator struct {
	numberCharacters int
}

func NewTrigramSeparator() Separator {
	return &trigramSeparator{
		numberCharacters: 3,
	}
}

func (t *trigramSeparator) Extract(text string) []string {
	words := strings.Split(text, " ")
	var grams []string
	for _, word := range words {
		grams = append(grams, t.extractWord(word)...)
	}
	return grams
}

func (t *trigramSeparator) extractWord(word string) []string {
	word = "  " + word + " "
	grams := make([]string, len(word)-t.numberCharacters+1)
	for i := 0; i < len(word)-t.numberCharacters+1; i++ {
		word := word[i : i+t.numberCharacters]
		grams[i] = word
	}
	return grams
}
