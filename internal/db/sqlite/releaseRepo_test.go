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

func TestUpdateReleaseCoverUrl(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		releaseRepo := NewReleaseSqliteRepo(db)
		artistRepo := NewArtistSqliteRepo(db)

		release := models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{},
		}
		artistId, _ := artistRepo.AddArtist("21 Savage")
		releaseRepo.AddRelease(release, artistId)

		err := releaseRepo.UpdateReleaseCoverUrl(release.Id, "https://cover.com")
		assert.NoError(t, err, "error didn't expected")

		releaseFromDB, _ := releaseRepo.GetReleaseById(release.Id)
		assert.Equal(t, releaseFromDB.CoverUrl, "https://cover.com")
	})
}

func TestGetReleasesByMonth(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)
		releases := []models.Release{
			{
				1,
				models.Artist{"21 Savage"},
				"American Dream",
				models.Album,
				types.NewCustomDate(2024, time.January, 12),
				models.CoverUrl{},
			},
			{
				2,
				models.Artist{"Drake"},
				"Another Release",
				models.Album,
				types.NewCustomDate(2024, time.January, 20),
				models.CoverUrl{},
			},
			{
				3,
				models.Artist{"Eminem"},
				"Some Release",
				models.Album,
				types.NewCustomDate(2024, time.March, 1),
				models.CoverUrl{},
			},
		}
		artistRepo := NewArtistSqliteRepo(db)
		releaseRepo := NewReleaseSqliteRepo(db)

		for _, release := range releases {
			artistId, _ := artistRepo.AddArtist(release.Artist.Name)
			releaseRepo.AddRelease(release, artistId)
		}

		got, err := releaseRepo.GetReleasesByMonth(time.January, 2024, 2, 0)
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(t, 2, len(got))
		assert.Equal(t, got[0].Title, "American Dream")
		assert.Equal(t, got[1].Title, "Another Release")
	})
}

func TestGetReleasesByYear(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)
		releases := []models.Release{
			{
				1,
				models.Artist{"21 Savage"},
				"American Dream",
				models.Album,
				types.NewCustomDate(2024, time.January, 12),
				models.CoverUrl{},
			},
			{
				2,
				models.Artist{"Drake"},
				"Another Release",
				models.Album,
				types.NewCustomDate(2024, time.January, 20),
				models.CoverUrl{},
			},
			{
				3,
				models.Artist{"Eminem"},
				"Some Release",
				models.Album,
				types.NewCustomDate(2023, time.March, 1),
				models.CoverUrl{},
			},
		}
		artistRepo := NewArtistSqliteRepo(db)
		releaseRepo := NewReleaseSqliteRepo(db)

		for _, release := range releases {
			artistId, _ := artistRepo.AddArtist(release.Artist.Name)
			releaseRepo.AddRelease(release, artistId)
		}

		got, err := releaseRepo.GetReleasesByYear(2024, 10, 0)
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(t, 2, len(got))
		assert.Equal(t, got[0].Title, "American Dream")
		assert.Equal(t, got[1].Title, "Another Release")
	})

	t.Run("not found releases in 2023", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)
		releases := []models.Release{
			{
				1,
				models.Artist{"21 Savage"},
				"American Dream",
				models.Album,
				types.NewCustomDate(2024, time.January, 12),
				models.CoverUrl{},
			},
			{
				2,
				models.Artist{"Drake"},
				"Another Release",
				models.Album,
				types.NewCustomDate(2024, time.January, 20),
				models.CoverUrl{},
			},
		}
		artistRepo := NewArtistSqliteRepo(db)
		releaseRepo := NewReleaseSqliteRepo(db)

		for _, release := range releases {
			artistId, _ := artistRepo.AddArtist(release.Artist.Name)
			releaseRepo.AddRelease(release, artistId)
		}

		got, err := releaseRepo.GetReleasesByYear(2023, 10, 0)
		assert.ErrorIs(t, err, ErrReleasesNotFound)
		assert.Nil(t, got)
	})
}

func TestGetReleasesWithoutCover(t *testing.T) {
	t.Run("we have cover in database", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)
		release := models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{},
		}

		releaseRepo := NewReleaseSqliteRepo(db)
		artistRepo := NewArtistSqliteRepo(db)

		artistId, _ := artistRepo.AddArtist(release.Artist.Name)
		releaseRepo.AddRelease(release, artistId)

		got, err := releaseRepo.GetReleasesWithoutCover()
		assert.NoError(t, err, "error didn't expected")
		assert.Equal(t, 1, len(got))
		assert.Equal(t, got[0].Title, "American Dream")
		assert.Equal(t, got[0].CoverUrl, "")
	})

	t.Run("all releases with cover", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)
		release := models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{
				Value:   "https://cover.com",
				IsValid: true,
			},
		}

		releaseRepo := NewReleaseSqliteRepo(db)
		artistRepo := NewArtistSqliteRepo(db)

		artistId, _ := artistRepo.AddArtist(release.Artist.Name)
		releaseRepo.AddRelease(release, artistId)

		got, err := releaseRepo.GetReleasesWithoutCover()
		assert.ErrorIs(t, err, ErrReleasesNotFound)
		assert.Nil(t, got)
	})
}

func TestGetReleasesByDay(t *testing.T) {
	t.Run("sucess case", func(t *testing.T) {
		db := prepareTestDb(t)
		defer removeTestDB(t, db)

		release := models.Release{
			1,
			models.Artist{"21 Savage"},
			"American Dream",
			models.Album,
			types.NewCustomDate(2024, time.January, 12),
			models.CoverUrl{
				Value:   "https://cover.com",
				IsValid: true,
			},
		}

		artistRepo := NewArtistSqliteRepo(db)
		releaseRepo := NewReleaseSqliteRepo(db)

		artistId, _ := artistRepo.AddArtist(release.Artist.Name)
		releaseRepo.AddRelease(release, artistId)

		got, err := releaseRepo.GetReleasesByDay(2024, time.January, 12, 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, len(got), 1)
		assert.Equal(t, "American Dream", got[0].Title)
	})
}
