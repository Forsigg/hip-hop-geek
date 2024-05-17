package sqlite

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/models"
)

type SqliteRepository struct {
	db.ArtistsRepositoryInterface
	db.ReleaseRepositoryInterface
	db.UsersRepositoryInterface
}

func NewSqliteRepository(db *sqlx.DB) *SqliteRepository {
	return &SqliteRepository{
		NewArtistSqliteRepo(db),
		NewReleaseSqliteRepo(db),
		NewUserSqliteRepo(db),
	}
}

func (s *SqliteRepository) Close() {
	s.CloseArtistRepo()
	s.CloseReleaseRepo()
}

func (s *SqliteRepository) CreateReleaseWithArtist(release models.Release) (int, error) {
	artistId := 0
	artist, err := s.GetArtistByName(release.Artist.Name)
	if err != nil {
		if strings.Contains(
			err.Error(),
			fmt.Sprintf("artist with name %s not found", release.Artist.Name),
		) {
			artistId, err = s.AddArtist(release.Artist.Name)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	if artist != nil {
		artistId = artist.Id
	}
	releaseId, err := s.AddRelease(release, artistId)
	if err != nil {
		return 0, err
	}

	return releaseId, nil
}

func (s *SqliteRepository) CreateMultiArtistsAndReleases(releases []models.Release) error {
	log.Println("start inserting releases into db...")
	artistIdToReleaseMap := make(map[int][]models.Release, 0)
	for _, release := range releases {
		artist, err := s.GetArtistByName(release.Artist.Name)
		if err != nil {
			if strings.Contains(
				err.Error(),
				fmt.Sprintf("artist with name %s not found", release.Artist.Name),
			) {
				log.Println(err.Error())
				artistId, err := s.AddArtist(release.Artist.Name)
				if err != nil {
					return err
				}
				artistIdToReleaseMap[artistId] = append(artistIdToReleaseMap[artistId], release)
			} else {
				return err
			}
		} else {
			artistIdToReleaseMap[artist.Id] = append(artistIdToReleaseMap[artist.Id], release)
		}
	}

	for artistId, releasesArr := range artistIdToReleaseMap {
		for _, release := range releasesArr {
			_, err := s.AddRelease(release, artistId)
			if err != nil {
				if errors.Is(err, ErrReleaseAlreadyExists) {
					continue
				}
				log.Printf("inserted release %s - %s", release.Artist.Name, release.Title)
				return err
			}
		}
	}

	return nil
}
