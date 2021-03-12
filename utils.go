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

func calcSimilarity(matchCount int, len1 int, len2 int) float32 {
	return float32(matchCount) / float32(len1+len2-matchCount)
}

func iterateWordSimilarity(gramIndexes2 []int, foundIndexes []bool, uniqueGrams1 int) float32 {
	var smlrMax float32
	var smlrCur float32

	lastPos := make([]int, len(foundIndexes))
	for i := range lastPos {
		lastPos[i] = -1
	}

	lower := -1
	uniqueGrams2 := 0
	count := 0
	upper := -1

	for i := 0; i < len(gramIndexes2); i++ {
		trgIndex := gramIndexes2[i]

		if lower >= 0 || foundIndexes[trgIndex] {
			if lastPos[trgIndex] < 0 {
				uniqueGrams2++
				if foundIndexes[trgIndex] {
					count++
				}

			}
			lastPos[trgIndex] = i
		}

		if foundIndexes[trgIndex] {
			upper = i
			if lower == -1 {
				lower = i
				uniqueGrams2 = 1
			}

			smlrCur = calcSimilarity(count, uniqueGrams1, uniqueGrams2)

			tmpCount := count
			tmpUniqueGram2 := uniqueGrams2
			prevLower := lower

			for tmpLower := lower; tmpLower <= upper; tmpLower++ {
				smlrTmp := calcSimilarity(tmpCount, uniqueGrams1, tmpUniqueGram2)
				if smlrTmp > smlrCur {
					smlrCur = smlrTmp
					uniqueGrams2 = tmpUniqueGram2
					lower = tmpLower
					count = tmpCount
				}

				tmpTrgIndex := gramIndexes2[tmpLower]
				if lastPos[tmpTrgIndex] == tmpLower {
					tmpUniqueGram2--
					if foundIndexes[tmpTrgIndex] {
						tmpCount--
					}
				}
			}

			smlrMax = float32(math.Max(float64(smlrMax), float64(smlrCur)))

			for tmpLower := prevLower; tmpLower < lower; tmpLower++ {

				tmpTrgIndex := gramIndexes2[tmpLower]
				if lastPos[tmpTrgIndex] == tmpLower {
					lastPos[tmpTrgIndex] = -1
				}
			}
		}
	}

	return smlrMax
}
