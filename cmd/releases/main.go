package main

import (
	"log"
	"time"

	"hip-hop-geek/internal/fetcher"
	"hip-hop-geek/internal/services/releases"
)

func main() {
	fetcher := fetcher.NewHipHopDXFetcher()
	service := releases.NewHipHopDXService(fetcher)
	releases := service.GetMonthReleases(2024, time.January, true)

	for _, rel := range releases {
		log.Println(rel)
	}
}
