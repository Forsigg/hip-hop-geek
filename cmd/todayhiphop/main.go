package main

import (
	"fmt"
	"log"

	"hip-hop-geek/internal/fetcher"
)

func main() {
	evenstFetcher := fetcher.NewTodayHipHopFetcher()
	defer evenstFetcher.Close()
	event, err := evenstFetcher.GetTodayEvent()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(event)
}
