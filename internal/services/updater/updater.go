package updater

import (
	"fmt"
	"log"
	"time"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/models"
	"hip-hop-geek/pkg/covers"
)

var allMonths = []time.Month{
	time.January, time.February, time.March,
	time.April, time.May, time.June,
	time.July, time.August, time.September,
	time.October, time.November, time.December,
}

type HipHopService interface {
	GetMonthReleases(year int, month time.Month, withCover bool) []models.Release
	GetAllYearReleases(year int, withCover bool) []models.Release
	GetAllYearSingles(year int, withCover bool) []models.Release
}

type ReleaseRepositoryInterface interface {
	AddRelease(release models.Release, artId int) (int, error)
	GetReleaseById(id int) (*db.ReleaseDB, error)
	GetReleaseByTitle(title string) (*db.ReleaseDB, error)
	UpdateReleaseCoverUrl(releaseId int, coverUrl string) error
}

type ArtistsRepositoryInterface interface {
	AddArtist(artistName string) (int, error)
	GetArtistByName(artistName string) (*db.ArtistDB, error)
	GetArtistById(id int) (*db.ArtistDB, error)
}

type DbRepository interface {
	ReleaseRepositoryInterface
	ArtistsRepositoryInterface
	CreateReleaseWithArtist(release models.Release) (int, error)
	CreateMultiArtistsAndReleases(releases []models.Release) error
}

type Updater struct {
	HipHopService
	DbRepository
}

func NewUpdater(hipHopService HipHopService, dbRepo DbRepository) *Updater {
	return &Updater{
		hipHopService,
		dbRepo,
	}
}

func (u *Updater) StartUploadReleases(timeToUpdate time.Duration, years []int, withCover bool) {
	ticker := time.NewTicker(timeToUpdate)

	for {
		select {
		case <-ticker.C:
			for _, year := range years {
				releases := u.GetAllYearReleases(year, withCover)
				err := u.CreateMultiArtistsAndReleases(releases)
				if err != nil {
					log.Fatal(err)
				}
			}
		default:
			log.Println("i will be wait 1 sec")
			time.Sleep(1 * time.Second)
		}
	}
}

func (u *Updater) CreateReleasesInDB(releases []models.Release) error {
	if err := u.CreateMultiArtistsAndReleases(releases); err != nil {
		return err
	}
	return nil
}

func (u *Updater) UpdateCoversInReleases(releases []models.Release) error {
	for _, release := range releases {
		query := fmt.Sprintf("%s - %s", release.Artist, release.Title)
		cover := covers.NewCoverBook().GetCoverByQuery(query, 500)
		u.UpdateReleaseCoverUrl(release.Id)
	}
}
