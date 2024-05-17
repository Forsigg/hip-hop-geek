package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"hip-hop-geek/internal/bot"
	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/fetcher"
	"hip-hop-geek/internal/services/releases"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./internal/db/db.db")
	sqliteRepo := sqlite.NewSqliteRepository(db)
	eventFetcher := fetcher.NewTodayHipHopFetcher()
	releaseFetcher := fetcher.NewHipHopDXFetcher()
	service := releases.NewHipHopService(sqliteRepo, releaseFetcher, eventFetcher)
	token := "6871130366:AAF2LkBSSJpvcRl0mqliDu5zTEqAfaiADZc"
	ctx, cancel := context.WithCancel(context.Background())
	bot := bot.NewTGBot(token, service)

	go bot.Start(ctx, 30)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("received signal to stop progamm...")
	cancel()
	log.Println("bot stopped")
}
