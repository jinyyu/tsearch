package tsearch

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
