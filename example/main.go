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
	storage.Speak()
}
