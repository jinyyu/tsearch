package main

import (
	"bufio"
	"github.com/jinyyu/tsearch"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Desc   string
	Pinyin string
}

func loadTestData(path string) (items map[uint32]*Item, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	items = make(map[uint32]*Item)

	reader := bufio.NewReader(strings.NewReader(string(data)))
	for {
		line, _, e := reader.ReadLine()
		if e != nil {
			break
		}

		result := strings.Split(string(line), ",")
		if len(result) != 4 {
			//log.Printf("invalid record %v", result)
			continue
		}

		idStr := result[0]
		desc := result[2]
		pingyin := result[3]
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			log.Fatalf("invalid record %v", err)
			continue
		}

		items[uint32(id)] = &Item{
			Desc:   desc,
			Pinyin: pingyin,
		}
	}
	return
}

func main() {
	storage, err := tsearch.NewRedisStorage("127.0.0.1:6379")
	if err != nil {
		log.Fatal("NewRedisStorage error")
	}

	separator := tsearch.NewTrigramSeparator()

	textSearch := tsearch.NewTextSearch(separator, storage)

	items, err := loadTestData("./test_data.csv")
	if err != nil {
		log.Fatalf("loadTestData error %v", err)
	}

	for id, item := range items {
		err = textSearch.UpdateText(id, item.Pinyin)
		if err != nil {
			log.Fatalf("UpdateText error %v", err)
		}
	}

	start := time.Now()
	log.Printf("search start %s", start.String())
	results, err := textSearch.Search("an quan niu lan qi", 0.6, 10)

	if err != nil {
		log.Fatalf("Search error %v", err)
	}
	end := time.Now()
	log.Printf("search end %s", end.String())

	for _, result := range results {
		item, _ := items[result.ID]
		log.Printf("id =%d, similarity =  %f, desc = %s", result.ID, result.Similarity, item.Desc)
	}

	log.Printf("use time %d (ms)", (end.UnixNano()-start.UnixNano())/1000000)
}
