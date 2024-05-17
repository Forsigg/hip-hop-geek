package sqlite

import (
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/models"
	"hip-hop-geek/internal/types"
)

func TestReleaseRepo(t *testing.T) {
	t.Run("check db creates", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		fs := os.DirFS(".")
		file, _ := fs.Open("test.db")
		fileStats, _ := file.Stat()
		if fileStats.Name() != "test.db" {
			t.Errorf("test.db not created")
		}
	})

	t.Run("check release creates correct", func(t *testing.T) {
		expected := &db.ReleaseDB{
			Id:       1,
			Artist:   db.ArtistDB{1, "21 Savage"},
			Title:    "American Dream",
			Type:     models.Album,
			OutYear:  2024,
			OutMonth: 1,
			OutDay:   12,
			CoverUrl: "",
		}
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		artistRepo := NewArtistSqliteRepo(db)
		artId, _ := artistRepo.AddArtist("21 Savage")
		repo := NewReleaseSqliteRepo(db)
		id, err := repo.AddRelease(models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{},
		}, artId)
		if err != nil {
			t.Errorf("error didn't expected: %s", err)
		}

		got, err := repo.GetReleaseById(id)
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(t, got, expected, "incorrect release added: want %v got %v", expected, got)
	})

	t.Run("check get release by id work correct", func(t *testing.T) {
		expected := &db.ReleaseDB{
			Id:       1,
			Artist:   db.ArtistDB{1, "21 Savage"},
			Title:    "American Dream",
			Type:     models.Album,
			OutYear:  2024,
			OutMonth: 1,
			OutDay:   12,
			CoverUrl: "",
		}
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		artistRepo := NewArtistSqliteRepo(db)
		artId, _ := artistRepo.AddArtist("21 Savage")
		repo := NewReleaseSqliteRepo(db)
		id, err := repo.AddRelease(models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{},
		}, artId)

		got, err := repo.GetReleaseById(id)
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(t, got, expected, "incorrect release get by id: want %v got %v", expected, got)
	})

	t.Run("check GetArtistByName work correct", func(t *testing.T) {
		expected := &db.ReleaseDB{
			Id:       1,
			Artist:   db.ArtistDB{1, "21 Savage"},
			Title:    "American Dream",
			OutYear:  2024,
			Type:     models.Album,
			OutMonth: 1,
			OutDay:   12,
			CoverUrl: "",
		}
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		artistRepo := NewArtistSqliteRepo(db)
		artId, _ := artistRepo.AddArtist("21 Savage")
		repo := NewReleaseSqliteRepo(db)
		relForLoad := models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{},
		}
		repo.AddRelease(relForLoad, artId)

		got, err := repo.GetReleaseByTitle(relForLoad.Title)
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(
			t,
			got,
			expected,
			"incorrect release get by name: want %v got %v",
			expected,
			got,
		)
	})
}
