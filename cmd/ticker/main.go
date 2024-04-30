package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/fetcher"
	"hip-hop-geek/internal/services/releases"
	"hip-hop-geek/internal/services/updater"
)

func main() {
	fetcher := fetcher.NewHipHopDXFetcher()
	service := releases.NewHipHopDXService(fetcher)
	db := sqlx.MustOpen("sqlite3", "./internal/db/db.db")
	repo := sqlite.NewSqliteRepository(db)
	updater := updater.NewUpdater(service, repo)

	ch := make(chan struct{})
	go updater.StartUploadReleases(10*time.Second, []int{2023, 2024}, true)

	x := <-ch
	log.Println(x)
}
