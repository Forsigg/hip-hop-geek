package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	service := releases.NewHipHopService(fetcher)
	db := sqlx.MustOpen("sqlite3", "./internal/db/db.db")
	repo := sqlite.NewSqliteRepository(db)
	updater := updater.NewUpdater(service, repo)

	sigCh := make(chan os.Signal, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go updater.StartUploadReleases(ctx, 10*time.Second, []int{2023, 2024}, false)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("cancelling...")
	cancel()
	log.Println("all closed, stop program")
	// x := <-ch
	// log.Println(x)
}
