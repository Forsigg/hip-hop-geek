package main

import (
	"log"

	"hip-hop-geek/pkg/covers"
)

func main() {
	cover := covers.NewCoverBook().GetCoverByQuery("Lil Baby - Freestyle", 500).Url
	log.Println(cover)
}
