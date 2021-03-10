package tsearch

import (
	"sort"
	"strings"
)

type TextSearch struct {
	separator Separator
	storage   Storage
}

func NewTextSearch(separator Separator, storage Storage) *TextSearch {
	return &TextSearch{
		separator: separator,
		storage:   storage,
	}
}

func (t *TextSearch) UpdateText(id uint32, text string) (err error) {
	oldTokens, err := t.storage.GetTokens(id)
	if err != nil {
		return err
	}
	newTokens := t.separator.Extract(text)

	err = t.storage.SaveTokens(id, newTokens)
	if err != nil {
		return err
	}

	return t.storage.UpdateIndex(id, oldTokens, newTokens)
}

// WordSimilarity 计算两个文本的相似度，主要用于调试
func (t *TextSearch) WordSimilarity(a string, b string) float32 {
	tokenA := t.separator.Extract(a)
	tokenA = DistinctTokens(tokenA)

	tokenB := DistinctTokens(t.separator.Extract(b))
	sort.Strings(tokenB)
	println(strings.Join(tokenB, ","))
	tokenMap := map[string]bool{}
	for _, tk := range tokenB {
		tokenMap[tk] = true
	}

	var matchCount int
	for _, tk := range tokenA {
		_, ok := tokenMap[tk]
		if ok {
			matchCount += 1
		}
	}
	return float32(matchCount) / float32(len(tokenA))
}

type SearchResult struct {
	ID         uint32
	Similarity float32
}

func (t *TextSearch) Search(text string, similarityThreshold float32, limit int) (results []*SearchResult, err error) {
	tokens := t.separator.Extract(text)
	if len(tokens) == 0 {
		return
	}

	counters, err := t.storage.SearchIndex(tokens)
	if err != nil {
		return
	}

	numberTokens := float32(len(tokens))
	for id, counter := range counters {
		similarity := float32(counter.Count) / numberTokens
		if similarity >= similarityThreshold {
			results = append(results, &SearchResult{
				ID:         id,
				Similarity: similarity,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if len(results) > limit {
		results = results[0:limit]
	}
	return
}
