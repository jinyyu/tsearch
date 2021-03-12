package main

import (
	"github.com/jinyyu/tsearch"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Item struct {
	ID     uint32 `db:"id"`
	Desc   string `db:"description"`
	Pinyin string `db:"ping_yin"`
}

func loadTestData() (result map[uint32]*Item, err error) {
	// this Pings the database trying to connect
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", "user=ljy dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var items []Item

	err = db.Select(&items, "select id,description,ping_yin from test_software")
	if err != nil {
		log.Fatalln(err)
	}
	result = map[uint32]*Item{}
	for i := range items {
		item := &items[i]
		result[item.ID] = item
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

	items, err := loadTestData()
	if err != nil {
		log.Fatalf("loadTestData error %v", err)
	}
	/*

		for id, item := range items {
			err = textSearch.UpdateText(id, item.Pinyin)
			if err != nil {
				log.Fatalf("UpdateText error %v", err)
			}
		}

	*/

	start := time.Now()
	log.Printf("search start %s", start.String())
	results, err := textSearch.Search("niu lan qi an quan", 0.6, 100)

	if err != nil {
		log.Fatalf("Search error %v", err)
	}
	end := time.Now()
	log.Printf("search end %s", end.String())

	for _, result := range results {
		item, _ := items[result.ID]
		log.Printf("id = %d, similarity = %f, desc = %s", result.ID, result.Similarity, item.Desc)
	}

	log.Printf("use time %d (ms)", (end.UnixNano()-start.UnixNano())/1000000)
}
