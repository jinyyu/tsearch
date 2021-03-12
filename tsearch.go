package tsearch

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
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

	err = t.dropGrams(id)
	if err != nil {
		return err
	}

	err = t.insertGrams(id, text)
	return err
}

func (t *TextSearch) dropGrams(id uint32) (err error) {
	idDistinctKey := t.getDistinctIDKey(id)
	idStr := fmt.Sprintf("%d", id)

	//获取旧的的分词
	values, err := t.storage.MultiGet(idDistinctKey)
	if err != nil {
		return err
	}
	if values[0] != "" {
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
	return nil
}

func (t *TextSearch) insertGrams(id uint32, text string) (err error) {
	idKey := t.getIDKey(id)
	idDistinctKey := t.getDistinctIDKey(id)

	//更新分词
	grams := t.separator.Extract(text)

	updateParam := make([]*KeyValue, 0, 2)
	if len(grams) > 0 {
		data, _ := json.Marshal(grams)
		updateParam = append(updateParam, &KeyValue{
			Key:   idKey,
			Value: string(data),
		})
	}

	distinctGrams := distinctStrings(grams)
	if len(distinctGrams) > 0 {
		data, _ := json.Marshal(distinctGrams)
		updateParam = append(updateParam, &KeyValue{
			Key:   idDistinctKey,
			Value: string(data),
		})
	}

	err = t.storage.MultiSet(updateParam...)
	if err != nil {
		return err
	}

	kvs := make([]*KeyValue, len(distinctGrams))
	for i := range distinctGrams {
		kvs[i] = &KeyValue{
			Key:   t.getGramKey(distinctGrams[i]),
			Value: fmt.Sprintf("%d", id),
		}
	}
	err = t.storage.MultiSetAdd(kvs...)
	if err != nil {
		return err
	}
	return err

}

func (t *TextSearch) getIDKey(id uint32) string {
	return fmt.Sprintf("id:%d", id)
}

func (t *TextSearch) getDistinctIDKey(id uint32) string {
	return fmt.Sprintf("id_distinct:%d", id)
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
	distinctGrams := distinctStrings(grams)

	hists, err := t.calcGramHits(distinctGrams, similarityThreshold)
	if err != nil {
		return nil, err
	}

	if len(hists) == 0 {
		return
	}

	idKeys := make([]string, 0, len(hists))
	for i := range hists {
		id := uint32(hists[i])
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

		s := calcWordSimilarity(grams, tokens)
		if s < similarityThreshold {
			continue
		}
		results = append(results, &SearchResult{
			ID:         uint32(hists[i]),
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

func (t *TextSearch) calcGramHits(grams []string, similarityThreshold float32) (ids []uint64, err error) {
	keys := make([]string, len(grams))
	for i := range grams {
		keys[i] = t.getGramKey(grams[i])
	}
	replies, err := t.storage.MultiGetMembers(keys...)
	if err != nil {
		return nil, err
	}

	idMaps := map[uint64]int{}

	for _, reply := range replies {
		hitIDs, err := redis.Uint64s(reply, nil)
		if err != nil {
			return nil, err
		}

		for _, id := range hitIDs {
			idMaps[id] += 1
		}
	}

	atLeast := int(float32(len(grams)) * similarityThreshold)
	for id, count := range idMaps {
		if count < atLeast {
			continue
		}

		ids = append(ids, id)
	}
	return
}
