package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"

	"hip-hop-geek/internal/bot"
	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/fetcher"
	"hip-hop-geek/internal/services/releases"
	"hip-hop-geek/internal/services/updater"
	"hip-hop-geek/internal/utils"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	baseProjDir := os.Getenv("BASE_PROJ_DIR")
	dbPath := filepath.Join(baseProjDir, "database", "db.db")
	logPath := filepath.Join(baseProjDir, "logs", "bot.log")
	migraionsDir := filepath.Join(baseProjDir, "internal", "db", "migrations")

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)
	log.Println("Logger activated")

	// prepare and check logger, db file and apply migrations
	if err = utils.EnsureFileExistsAndCreate(dbPath); err != nil {
		log.Fatal(err)
	}

	if err = migrationDBUp(dbPath, migraionsDir); err != nil {
		log.Fatal(err)
	}

	// init database and repository layer
	db := sqlx.MustConnect("sqlite3", dbPath)
	repo := sqlite.NewSqliteRepository(db)

	// init event and release fetchers layer
	eventFetcher := fetcher.NewTodayHipHopFetcher()
	releaseFetcher := fetcher.NewHipHopDXFetcher()

	// init service layer
	service := releases.NewHipHopService(repo, releaseFetcher, eventFetcher)

	// init updater
	updater := updater.NewUpdater(service, repo)

	// prepare context
	ctx, cancel := context.WithCancel(context.Background())

	// init bot
	bot := bot.NewTGBot(os.Getenv("TG_BOT_TOKEN"), service, updater)
	if bot == nil {
		log.Fatal("bot nil, the end")
	}

	timeForUpdate := time.Duration(8 * time.Hour)
	// start goroutines with update releases and tg-bot
	go updater.StartUploadReleases(ctx, timeForUpdate, []int{2023, 2024}, false)
	go bot.Start(ctx, 30)
	go bot.SendEventAndReleasesEveryday(ctx)

	// chan for os signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// blocking operation - wait for signal
	<-sigCh
	log.Println("received signal to stop progamm...")
	cancel()
	log.Println("bot stopped")
	logFile.Close()
}

// migrationDBUp выполняет миграции базы данных SQLite с использованием goose.
//
// path - путь к файлу базы данных SQLite.
// migrationsPath - путь к директории с файлами миграций.
//
// Функция выполняет следующие действия:
// 1. Открывает соединение с базой данных SQLite по указанному пути.
// 2. Устанавливает диалект базы данных для goose как "sqlite3".
// 3. Применяет все миграции, найденные в указанной директории миграций.
//
// В случае возникновения ошибки при открытии соединения с базой данных или
// при применении миграций, функция возвращает ошибку, обернутую с дополнительным
// сообщением для облегчения диагностики проблемы.
func migrationDBUp(path string, migrationsPath string) error {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("error while connecting to test db: %w", err)
	}

	goose.SetDialect("sqlite3")
	err = goose.Up(db.DB, migrationsPath)
	if err != nil {
		return fmt.Errorf("error while applying migrations: %w", err)
	}

	return nil
}
