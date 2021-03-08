package main

import (
	"github.com/jinyyu/tsearch"
	"log"
)

func main() {
	storage, err := tsearch.NewRedisStorage("127.0.0.1:6379")
	if err != nil {
		log.Fatal("NewRedisStorage error")
	}

	separator := tsearch.NewTrigramSeparator()

	textSearch := tsearch.NewTextSearch(separator, storage)
	err = textSearch.UpdateText(1, "two words")
	if err != nil {
		log.Fatalf("UpdateText error %v", err)
	}

	results, err := textSearch.Search("word", 0.8)
	if err != nil {
		log.Fatalf("Search error %v", err)
	}

	for _, result := range results {
		log.Printf("id = %d, similarity = %f", result.ID, result.Similarity)
	}
}
