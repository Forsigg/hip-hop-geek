package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"hip-hop-geek/internal/db/sqlite"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "./internal/db/db.db")
	repo := sqlite.NewSqliteRepository(db)
	releaseDb, err := repo.GetReleaseByTitle("American Dream")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(releaseDb)
}
