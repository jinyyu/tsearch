package tsearch

import (
	"encoding/json"
	"fmt"
	"sort"
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
	idStr := fmt.Sprintf("%d", id)
	idKey := t.getIDKey(id)
	//获取旧的的分词
	values, err := t.storage.MultiGet(idKey)
	if err != nil {
		return err
	}
	if len(values) > 0 && values[0] != "" {
		var oldGrams []string
		err = json.Unmarshal([]byte(values[0]), &oldGrams)
		if err != nil {
			return err
		}

		kvs := make([]*KeyValue, len(oldGrams))
		for i, gram := range oldGrams {
			kvs[i] = &KeyValue{
				Key:   t.getGramKey(gram),
				Value: idStr,
			}
		}
		//删除旧的分词
		err = t.storage.MultiSetDel(kvs...)
		if err != nil {
			return err
		}
	}

	//更新分词
	grams := t.separator.Extract(text)
	data, err := json.Marshal(grams)
	if err != nil {
		return err
	}
	err = t.storage.MultiSet(&KeyValue{
		Key:   idKey,
		Value: string(data),
	})
	if err != nil {
		return err
	}

	kvs := make([]*KeyValue, len(grams))
	for i, gram := range grams {
		kvs[i] = &KeyValue{
			Key:   t.getGramKey(gram),
			Value: idStr,
		}
	}

	err = t.storage.MultiSetAdd(kvs...)
	return err
}

func (t *TextSearch) getIDKey(id uint32) string {
	return fmt.Sprintf("id:%d", id)
}

func (t *TextSearch) getGramKey(gram string) string {
	return fmt.Sprintf("gram:%s", gram)
}

// WordSimilarity 计算两个文本的相似度，主要用于调试
func (t *TextSearch) WordSimilarity(word1 string, word2 string) float32 {
	token1 := t.separator.Extract(word1)
	token2 := t.separator.Extract(word2)
	return calcWordSimilarity(token1, token2)
}

type SearchResult struct {
	ID         uint32
	Similarity float32
}

func (t *TextSearch) Search(text string, similarityThreshold float32, limit int) (results []*SearchResult, err error) {

	grams := t.separator.Extract(text)
	if len(grams) == 0 {
		return nil, nil
	}

	hists, err := t.calcGramHits(grams)
	if err != nil {
		return nil, err
	}

	var idKeys []string
	var idList []uint32

	for id, _ := range hists {
		idList = append(idList, id)
		idKeys = append(idKeys, t.getIDKey(id))
	}

	values, err := t.storage.MultiGet(idKeys...)
	if err != nil {
		return nil, err
	}

	for i, value := range values {
		var tokens []string
		err = json.Unmarshal([]byte(value), &tokens)
		if err != nil {
			return nil, err
		}
		count, _ := hists[idList[i]]
		if CALCSML(len(grams), count.Count, len(grams)) < similarityThreshold {
			continue
		}

		s := calcWordSimilarity(grams, tokens)
		if s < similarityThreshold {
			continue
		}
		results = append(results, &SearchResult{
			ID:         idList[i],
			Similarity: s,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if len(results) > limit {
		results = results[0:limit]
	}
	return
}

type HitCounter struct {
	Count int
}

func (t *TextSearch) calcGramHits(grams []string) (hits map[uint32]*HitCounter, err error) {
	keys := make([]string, len(grams))
	for i := range grams {
		keys[i] = t.getGramKey(grams[i])
	}

	memberMap, err := t.storage.MultiGetMembers(keys...)
	if err != nil {
		return nil, err
	}

	hits = make(map[uint32]*HitCounter)

	for _, members := range memberMap {
		for _, id := range members {

			counter, ok := hits[uint32(id)]
			if ok {
				counter.Count += 1
			} else {
				counter = &HitCounter{
					Count: 1,
				}
				hits[uint32(id)] = counter
			}
		}
	}
	return hits, nil
}
