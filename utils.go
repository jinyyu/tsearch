package tsearch

import (
	"math"
	"sort"
)

func distinctStrings(a []string) []string {
	sort.Strings(a)
	ret := make([]string, 0, len(a))
	for i := range a {
		if i > 0 {
			if a[i] != a[i-1] {
				ret = append(ret, a[i])
			}
		} else {
			ret = append(ret, a[i])
		}
	}
	return ret
}

// 带位置信息的分词
type positionGram struct {
	gram  string //分词
	index int    //分词的位置
}

// 计算两组分词的相似度
func calcWordSimilarity(gram1 []string, gram2 []string) float32 {
	positionGrams := makePositionalGrams(gram1, gram2)

	gram2Indexes := make([]int, len(gram2))
	foundIndexes := make([]bool, len(positionGrams))
	uniqueGrams1 := 0

	j := 0
	for i := 0; i < len(positionGrams); i++ {
		if i > 0 {
			changed := positionGrams[i-1].gram != positionGrams[i].gram
			if changed {
				if foundIndexes[j] {
					uniqueGrams1++
				}
				j++
			}

		}
		if positionGrams[i].index >= 0 {
			gram2Indexes[positionGrams[i].index] = j
		} else {
			foundIndexes[j] = true
		}
	}

	if foundIndexes[j] {
		uniqueGrams1++
	}

	// Run iterative procedure to find maximum similarity with word
	return iterateWordSimilarity(gram2Indexes, foundIndexes, uniqueGrams1)
}

func makePositionalGrams(grams1 []string, grams2 []string) (result []positionGram) {
	result = make([]positionGram, len(grams1)+len(grams2))
	for i := 0; i < len(grams1); i++ {
		result[i].gram = grams1[i]
		result[i].index = -1
	}

	for i := 0; i < len(grams2); i++ {
		j := i + len(grams1)
		result[j].gram = grams2[i]
		result[j].index = i
	}

	sort.Slice(result, func(i, j int) bool {
		p1 := &result[i]
		p2 := &result[j]

		if p1.gram == p2.gram {
			return p1.index < p2.index
		} else {
			return p1.gram < p2.gram
		}
	})
	return result
}

func CALCSML(count int, len1 int, len2 int) float32 {
	return float32(count) / float32(len1+len2-count)
}

func iterateWordSimilarity(token2Indexes []int, found []bool, numberUniqueToken1 int) float32 {
	var smlr_max float32
	var smlr_cur float32

	lastpos := make([]int, len(found))
	for i := range lastpos {
		lastpos[i] = -1
	}

	lower := -1
	ulen2 := 0
	count := 0
	upper := -1

	for i := 0; i < len(token2Indexes); i++ {
		trgindex := token2Indexes[i]

		if lower >= 0 || found[trgindex] {
			if lastpos[trgindex] < 0 {
				ulen2++
				if found[trgindex] {
					count++
				}

			}
			lastpos[trgindex] = i
		}

		if found[trgindex] {
			upper = i
			if lower == -1 {
				lower = i
				ulen2 = 1
			}

			smlr_cur = CALCSML(count, numberUniqueToken1, ulen2)

			tmp_count := count
			tmp_ulen2 := ulen2
			prev_lower := lower

			for tmp_lower := lower; tmp_lower <= upper; tmp_lower++ {
				if true {
					smlr_tmp := CALCSML(tmp_count, numberUniqueToken1, tmp_ulen2)
					if smlr_tmp > smlr_cur {
						smlr_cur = smlr_tmp
						ulen2 = tmp_ulen2
						lower = tmp_lower
						count = tmp_count
					}
				}
				tmp_trgindex := token2Indexes[tmp_lower]
				if lastpos[tmp_trgindex] == tmp_lower {
					tmp_ulen2--
					if found[tmp_trgindex] {
						tmp_count--
					}
				}
			}

			smlr_max = float32(math.Max(float64(smlr_max), float64(smlr_cur)))

			for tmp_lower := prev_lower; tmp_lower < lower; tmp_lower++ {

				tmp_trgindex := token2Indexes[tmp_lower]
				if lastpos[tmp_trgindex] == tmp_lower {
					lastpos[tmp_trgindex] = -1
				}
			}
		}
	}

	return smlr_max
}
