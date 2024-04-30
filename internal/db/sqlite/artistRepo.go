package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"hip-hop-geek/internal/db"
)

var _ db.ArtistsRepositoryInterface = (*ArtistSqliteRepo)(nil)

const (
	createArtistStmt     = `INSERT INTO artists (name) VALUES(?);`
	getArtistByIdQuery   = `SELECT * FROM artists WHERE artist_id = ?;`
	getArtistByNameQuery = `SELECT * FROM artists WHERE name = ?;`
)

type ArtistSqlite struct {
	Id   int    `db:"artist_id"`
	Name string `db:"name"`
}

type ArtistSqliteRepo struct {
	DB *sqlx.DB
}

func NewArtistSqliteRepo(db *sqlx.DB) *ArtistSqliteRepo {
	return &ArtistSqliteRepo{db}
}

func (a *ArtistSqliteRepo) AddArtist(artistName string) (int, error) {
	res, err := a.DB.Exec(createArtistStmt, artistName)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: artists.name") {
			return 0, ErrArtistAlreadyExists
		}
		return 0, fmt.Errorf("db error add artist: %s", err)
	}
	artId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error while inserting artist: %s", err)
	}

	return int(artId), nil
}

func (a *ArtistSqliteRepo) GetArtistByName(artistName string) (*db.ArtistDB, error) {
	var artist []ArtistSqlite
	err := a.DB.Select(&artist, getArtistByNameQuery, artistName)
	if err != nil {
		return nil, fmt.Errorf("error while querying artist by name: %s", err)
	}

	if len(artist) > 1 {
		return nil, errors.New("getting artist by id returns few values, but expected one")
	} else if len(artist) == 0 {
		return nil, fmt.Errorf("artist with name %s not found", artistName)
	}

	return &db.ArtistDB{
		Id:   artist[0].Id,
		Name: artist[0].Name,
	}, nil
}

func (a *ArtistSqliteRepo) GetArtistById(id int) (*db.ArtistDB, error) {
	var artist []ArtistSqlite
	err := a.DB.Select(&artist, getArtistByIdQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error while querying artist by id: %s", err)
	}

	if len(artist) > 1 {
		return nil, errors.New("getting artist by id returns few values, but expected one")
	} else if len(artist) == 0 {
		return nil, fmt.Errorf("artist with id %d not found", id)
	}

	return &db.ArtistDB{
		Id:   artist[0].Id,
		Name: artist[0].Name,
	}, nil
}

func (a *ArtistSqliteRepo) Close() {
	a.DB.Close()
}
