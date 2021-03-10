package tsearch

import "sort"

func DistinctTokens(tokens []string) []string {
	ret := make([]string, 0, len(tokens))
	m := map[string]bool{}
	for _, token := range tokens {
		_, ok := m[token]
		if ok {
			continue
		}
		m[token] = true
		ret = append(ret, token)
	}
	return ret
}

/* Trigram with position */
type positionGram struct {
	token string
	index int
}

func makePositionalGram(token1 []string, token2 []string) (result []positionGram) {
	result = make([]positionGram, len(token1)+len(token2))
	for i := 0; i < len(token1); i++ {
		result[i].token = token1[i]
		result[i].index = -1
	}

	for i := 0; i < len(token2); i++ {
		j := i + len(token1)
		result[j].token = token2[j]
		result[j].index = i
	}

	sort.Slice(result, func(i, j int) bool {
		p1 := &result[i]
		p2 := &result[j]

		if p1.token == p2.token {
			return p1.index < p2.index
		} else {
			return p1.token < p2.token
		}
	})
	return result
}

func calcWordSimilarity(token1 []string, token2 []string) float32 {
	pGrams := makePositionalGram(token1, token2)
	//TOTO
}
