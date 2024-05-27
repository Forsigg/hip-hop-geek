package sqlite

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"

	"hip-hop-geek/internal/db"
)

func TestArtistRepo(t *testing.T) {
	t.Run("check db creates", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		fs := os.DirFS(".")
		file, err := fs.Open("test.db")
		assert.NoError(t, err, "open error: %s", err)

		fileStats, err := file.Stat()
		assert.NoError(t, err, "stat error: %s", err)
		if fileStats.Name() != "test.db" {
			t.Errorf("test.db not created")
		}
	})

	t.Run("check artist creates correct", func(t *testing.T) {
		expected := &db.ArtistDB{Id: 1, Name: "Lil Yachty"}
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewArtistSqliteRepo(db)
		id, err := repo.AddArtist("Lil Yachty")
		if err != nil {
			t.Errorf("error didn't expected: %s", err)
		}

		got, err := repo.GetArtistById(id)
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(t, got, expected, "incorrect artist added: want %v got %v", expected, got)
	})

	t.Run("check get by id work correct", func(t *testing.T) {
		expected := []db.ArtistDB{
			{1, "21 Savage"},
			{2, "Lil Yachty"},
			{3, "Drake"},
		}

		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewArtistSqliteRepo(db)
		idOne, _ := repo.AddArtist("21 Savage")
		idTwo, _ := repo.AddArtist("Lil Yachty")
		idThree, _ := repo.AddArtist("Drake")

		artOne, _ := repo.GetArtistById(idOne)
		artTwo, _ := repo.GetArtistById(idTwo)
		artThree, _ := repo.GetArtistById(idThree)

		assert.Equal(t, &expected[0], artOne,
			"artist not equals: want %v got %v", &expected[0], artOne,
		)
		assert.Equal(t, &expected[1], artTwo,
			"artist not equals: want %v got %v", &expected[1], artTwo,
		)
		assert.Equal(t, &expected[2], artThree,
			"artist not equals: want %v got %v", &expected[2], artThree,
		)
	})

	t.Run("check GetArtistByName work correct", func(t *testing.T) {
		expected := []db.ArtistDB{
			{1, "21 Savage"},
			{2, "Lil Yachty"},
			{3, "Drake"},
		}

		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		repo := NewArtistSqliteRepo(db)
		repo.AddArtist("21 Savage")
		repo.AddArtist("Lil Yachty")
		repo.AddArtist("Drake")

		artOne, _ := repo.GetArtistByName("21 Savage")
		artTwo, _ := repo.GetArtistByName("Lil Yachty")
		artThree, _ := repo.GetArtistByName("Drake")

		assert.Equal(t, &expected[0], artOne,
			"artist not equals: want %v got %v", &expected[0], artOne,
		)
		assert.Equal(t, &expected[1], artTwo,
			"artist not equals: want %v got %v", &expected[1], artTwo,
		)
		assert.Equal(t, &expected[2], artThree,
			"artist not equals: want %v got %v", &expected[2], artThree,
		)
	})
}

func prepareTestDb(t testing.TB) *sqlx.DB {
	t.Helper()

	testDbPath := "./test.db"
	migrationsDir := "../migrations"

	_, err := os.Create(testDbPath)
	if err != nil {
		log.Fatalf("error while creating test db: %s", err)
	}

	db, err := sqlx.Open("sqlite3", testDbPath)
	if err != nil {
		log.Fatalf("error while connecting to test db: %s", err)
	}

	// turn off migrations log, but catch error
	log.SetOutput(io.Discard)
	goose.SetDialect("sqlite3")
	err = goose.Up(db.DB, migrationsDir)
	// turn on logs in stdout
	log.SetOutput(os.Stdout)
	if err != nil {
		log.Fatalf("error while applying migrations: %s", err)
	}

	return db
}

func removeTestDB(t testing.TB, db *sqlx.DB) {
	t.Helper()
	db.Close()
	err := os.Remove("./test.db")
	if err != nil {
		log.Fatalf("error while deleting test database file: %s", err)
	}
}
