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
	var tokens []string
	for _, word := range words {
		tokens = append(tokens, t.extractWord(word)...)
	}
	return tokens
}

func (t *trigramSeparator) extractWord(word string) []string {
	word = "  " + word + " "
	tokens := make([]string, len(word)-t.numberCharacters+1)
	for i := 0; i < len(word)-t.numberCharacters+1; i++ {
		word := word[i : i+t.numberCharacters]
		tokens[i] = word
	}
	return tokens
}
