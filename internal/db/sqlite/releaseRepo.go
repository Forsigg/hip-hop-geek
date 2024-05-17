package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/models"
)

const (
	createReleaseStmt = `INSERT INTO releases
                        (release_id, artist_id, release_type, title, out_year, out_month, out_day, cover_url)
                        VALUES
                        (?, ?, ?, ?, ?, ?, ?, ?);`

	getReleaseByIdQuery = `
    SELECT r.release_id, a.artist_id AS "artist.artist_id", a.name AS "artist.name", r.title, r.out_year, r.out_month, r.out_day, r.cover_url, r.release_type
    FROM releases AS r
    JOIN artists AS a ON r.artist_id = "artist.artist_id" 
    WHERE r.release_id = ?;`

	getReleasesByMonthQuery = `
    SELECT r.release_id, a.artist_id AS "artist.artist_id", a.name AS "artist.name", r.title, r.out_year, r.out_month, r.out_day, r.cover_url, r.release_type
    FROM releases AS r
    JOIN artists AS a ON r.artist_id = "artist.artist_id" 
    WHERE r.out_month = ? AND r.out_year = ?
    ORDER BY r.out_year, r.out_month, r.out_day
    LIMIT ? OFFSET ?;`

	getReleasesWithoutCoverQuery = `
    SELECT r.release_id, a.artist_id AS "artist.artist_id", a.name AS "artist.name", r.title, r.out_year, r.out_month, r.out_day, r.cover_url, r.release_type
    FROM releases AS r
    JOIN artists AS a ON r.artist_id = "artist.artist_id" 
    WHERE r.cover_url = ""
    ORDER BY r.out_year, r.out_month, r.out_day;`

	getReleasesByYearQuery = `
    SELECT r.release_id, a.artist_id AS "artist.artist_id", a.name AS "artist.name", r.title, r.out_year, r.out_month, r.out_day, r.cover_url, r.release_type
    FROM releases AS r
    JOIN artists AS a ON r.artist_id = "artist.artist_id" 
    WHERE r.out_year = ?
    ORDER BY r.out_year, r.out_month, r.out_day
    LIMIT ? OFFSET ?;`

	getReleasesByNameQuery = `
    SELECT r.release_id, a.artist_id AS "artist.artist_id", a.name AS "artist.name", r.title, r.out_year, r.out_month, r.out_day, r.cover_url, r.release_type
    FROM releases AS r
    JOIN artists AS a ON r.artist_id = "artist.artist_id" 
    WHERE r.title = ?;`

	getReleasesByDayQuery = `
    SELECT r.release_id, a.artist_id AS "artist.artist_id", a.name AS "artist.name", r.title, r.out_year, r.out_month, r.out_day, r.cover_url, r.release_type
    FROM releases AS r
    JOIN artists AS a ON r.artist_id = "artist.artist_id" 
    WHERE r.out_year = ? AND r.out_month = ? AND r.out_day = ?
    LIMIT ? OFFSET ?;`

	updateReleaseCoverStmt = `
    UPDATE releases
    SET cover_url = ?
    WHERE release_id = ?;
    `
)

type ReleaseSqlite struct {
	Id       int            `db:"release_id"`
	Artist   ArtistSqlite   `db:"artist"`
	Title    string         `db:"title"`
	Type     int            `db:"release_type"`
	OutYear  int            `db:"out_year"`
	OutMonth int            `db:"out_month"`
	OutDay   int            `db:"out_day"`
	CoverUrl sql.NullString `db:"cover_url"`
}

type ReleaseSqliteRepo struct {
	DB *sqlx.DB
}

func NewReleaseSqliteRepo(db *sqlx.DB) *ReleaseSqliteRepo {
	return &ReleaseSqliteRepo{db}
}

func (r *ReleaseSqliteRepo) CloseReleaseRepo() {
	r.DB.Close()
}

func (r *ReleaseSqliteRepo) AddRelease(release models.Release, artId int) (int, error) {
	res, err := r.DB.Exec(
		createReleaseStmt,
		release.Id,
		artId,
		release.Type,
		release.Title,
		release.OutDate.Year(),
		release.OutDate.Month(),
		release.OutDate.Day(),
		release.CoverUrl.Value,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: releases.release_id") {
			return 0, ErrReleaseAlreadyExists
		}
		return 0, fmt.Errorf("db error add release: %s", err)
	}
	releaseId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error while inserting release: %s", err)
	}

	return int(releaseId), nil
}

func (r *ReleaseSqliteRepo) GetReleaseById(id int) (*db.ReleaseDB, error) {
	var releases []ReleaseSqlite
	err := r.DB.Select(&releases, getReleaseByIdQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error while querying release by id: %s", err)
	}

	if len(releases) > 1 {
		return nil, errors.New("getting release by id returns few values, but expected one")
	} else if len(releases) == 0 {
		return nil, fmt.Errorf("release with id %d not found", id)
	}

	return &db.ReleaseDB{
		Id:       releases[0].Id,
		Artist:   db.ArtistDB(releases[0].Artist),
		Title:    releases[0].Title,
		Type:     releases[0].Type,
		OutYear:  releases[0].OutYear,
		OutMonth: releases[0].OutMonth,
		OutDay:   releases[0].OutDay,
		CoverUrl: releases[0].CoverUrl.String,
	}, nil
}

func (r *ReleaseSqliteRepo) GetReleaseByTitle(title string) (*db.ReleaseDB, error) {
	var releases []ReleaseSqlite
	err := r.DB.Select(&releases, getReleasesByNameQuery, title)
	if err != nil {
		return nil, fmt.Errorf("error while querying release by name: %s", err)
	}

	if len(releases) > 1 {
		return nil, errors.New("getting release by name returns few values, but expected one")
	} else if len(releases) == 0 {
		return nil, fmt.Errorf("release with name %s not found", title)
	}

	return &db.ReleaseDB{
		Id:       releases[0].Id,
		Artist:   db.ArtistDB(releases[0].Artist),
		Title:    releases[0].Title,
		Type:     releases[0].Type,
		OutYear:  releases[0].OutYear,
		OutMonth: releases[0].OutMonth,
		OutDay:   releases[0].OutDay,
		CoverUrl: releases[0].CoverUrl.String,
	}, nil
}

func (r *ReleaseSqliteRepo) UpdateReleaseCoverUrl(releaseId int, coverUrl string) error {
	_, err := r.DB.Exec(updateReleaseCoverStmt, coverUrl, releaseId)
	if err != nil {
		return fmt.Errorf("error while updating cover_url in release(id %d): %w", releaseId, err)
	}
	return nil
}

func (r *ReleaseSqliteRepo) GetReleasesByMonth(
	month time.Month,
	year, limit, offset int,
) ([]*db.ReleaseDB, error) {
	var releasesFromDB []ReleaseSqlite

	var err error
	if month == 0 {
		err = r.DB.Select(&releasesFromDB, getReleasesByYearQuery, year, limit, offset)
	} else {
		err = r.DB.Select(&releasesFromDB, getReleasesByMonthQuery, month, year, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("error while querying release by month: %w", err)
	}

	if len(releasesFromDB) == 0 {
		return nil, ErrReleasesNotFound
	}

	releasesResult := make([]*db.ReleaseDB, 0, len(releasesFromDB))

	for _, rel := range releasesFromDB {
		releasesResult = append(releasesResult, &db.ReleaseDB{
			Id:       rel.Id,
			Artist:   db.ArtistDB(rel.Artist),
			Type:     rel.Type,
			Title:    rel.Title,
			OutYear:  rel.OutYear,
			OutMonth: rel.OutMonth,
			OutDay:   rel.OutDay,
			CoverUrl: rel.CoverUrl.String,
		})
	}

	return releasesResult, nil
}

func (r *ReleaseSqliteRepo) GetReleasesByYear(year, limit, offset int) ([]*db.ReleaseDB, error) {
	releases, err := r.GetReleasesByMonth(0, year, limit, offset)
	if err != nil {
		if errors.Is(err, ErrReleasesNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("error while querying release by year: %w", err)
	}

	return releases, nil
}

func (r *ReleaseSqliteRepo) GetReleasesWithoutCover() ([]*db.ReleaseDB, error) {
	var releasesFromDB []ReleaseSqlite

	err := r.DB.Select(&releasesFromDB, getReleasesWithoutCoverQuery)
	if err != nil {
		return nil, fmt.Errorf("error while getting releases without cover: %w", err)
	}

	if len(releasesFromDB) == 0 {
		return nil, ErrReleasesNotFound
	}

	releasesResult := make([]*db.ReleaseDB, 0, len(releasesFromDB))
	for _, rel := range releasesFromDB {
		releasesResult = append(releasesResult, &db.ReleaseDB{
			Id:       rel.Id,
			Artist:   db.ArtistDB(rel.Artist),
			Title:    rel.Title,
			Type:     rel.Type,
			OutYear:  rel.OutYear,
			OutMonth: rel.OutMonth,
			OutDay:   rel.OutDay,
			CoverUrl: rel.CoverUrl.String,
		})
	}

	return releasesResult, nil
}

func (r *ReleaseSqliteRepo) GetReleasesByDay(
	year int,
	month time.Month,
	day,
	limit,
	offset int,
) ([]*db.ReleaseDB, error) {
	var releasesFromDB []ReleaseSqlite
	err := r.DB.Select(&releasesFromDB, getReleasesByDayQuery, year, month, day, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while getting releases by day: %w", err)
	}

	if len(releasesFromDB) == 0 {
		return nil, ErrReleasesNotFound
	}

	releasesResult := make([]*db.ReleaseDB, 0, len(releasesFromDB))
	for _, rel := range releasesFromDB {
		releasesResult = append(releasesResult, &db.ReleaseDB{
			Id:       rel.Id,
			Artist:   db.ArtistDB(rel.Artist),
			Title:    rel.Title,
			Type:     rel.Type,
			OutYear:  rel.OutYear,
			OutMonth: rel.OutMonth,
			OutDay:   rel.OutDay,
			CoverUrl: rel.CoverUrl.String,
		})
	}

	return releasesResult, nil
}
