package updater

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/models"
	"hip-hop-geek/pkg/covers"
)

type HipHopService interface {
	FetchReleases(year int) ([]models.Release, error)
	FetchSingles(year int) ([]models.Release, error)

	GetMonthReleases(year int, month time.Month, limit, offset int) []models.Release
	GetAllYearReleases(year int, limit, offset int) []models.Release
	GetAllYearSingles(year int, withCover bool) []models.Release

	Close()
}

type ReleaseRepositoryInterface interface {
	AddRelease(release models.Release, artId int) (int, error)
	GetReleaseById(id int) (*db.ReleaseDB, error)
	GetReleaseByTitle(title string) (*db.ReleaseDB, error)
	GetReleasesWithoutCover() ([]*db.ReleaseDB, error)
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
	UpdateReleaseCoverUrl(releaseId int, coverUrl string) error
	Close()
}

type Updater struct {
	mu sync.Mutex
	HipHopService
	DbRepository
}

func NewUpdater(hipHopService HipHopService, dbRepo DbRepository) *Updater {
	return &Updater{
		sync.Mutex{},
		hipHopService,
		dbRepo,
	}
}

func (u *Updater) StartUploadReleases(
	ctx context.Context,
	timeToUpdate time.Duration,
	years []int,
	withCover bool,
) {
	log.Println("start updater on timer")
	ticker := time.NewTicker(timeToUpdate)

	// sleep at first run
	time.Sleep(30 * time.Second)

	u.RefreshReleases(years)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			u.RefreshReleases(years)
		case <-ctx.Done():
			u.Close()
		}
	}
}

func (u *Updater) Close() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		u.HipHopService.Close()
	}()
	go func() {
		defer wg.Done()
		u.DbRepository.Close()
	}()

	wg.Wait()
}

func (u *Updater) CreateReleasesInDB(releases []models.Release) error {
	if err := u.CreateMultiArtistsAndReleases(releases); err != nil {
		return err
	}
	return nil
}

func (u *Updater) UpdateCoversInDB() error {
	log.Println("start updating covers...")
	releases, err := u.GetReleasesWithoutCover()
	if err != nil {
		if errors.Is(sqlite.ErrReleasesNotFound, err) {
			log.Println("all releases with covers, cool")
			return nil
		}
		return err
	}

	var wg sync.WaitGroup
	for _, release := range releases {
		wg.Add(1)
		go func(release *db.ReleaseDB) {
			defer wg.Done()
			query := fmt.Sprintf("%s - %s", release.Artist.Name, release.Title)
			cover := covers.NewCoverBook().GetCoverByQuery(query, 600)
			if cover.Valid {
				log.Printf("set new cover for %s", query)
				err = u.UpdateReleaseCoverUrl(release.Id, cover.Url)
				if err != nil {
					log.Printf("error while updating cover for %s: %s", query, err)
				}
			}
		}(release)
	}

	wg.Wait()
	log.Println("all covers updated")
	return nil
}

func (u *Updater) RefreshReleases(years []int) {
	log.Println("looking for new releases")
	for _, year := range years {

		allReleases := make([]models.Release, 0, 10)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			log.Println("try to get new releases...")
			defer wg.Done()
			newReleases, err := u.FetchReleases(year)
			if err != nil {
				log.Fatal(err)
			}
			u.mu.Lock()
			allReleases = append(allReleases, newReleases...)
			u.mu.Unlock()
		}()

		go func() {
			log.Println("try to get new singles...")
			defer wg.Done()
			newSingles, err := u.FetchSingles(year)
			if err != nil {
				log.Fatal(err)
			}

			u.mu.Lock()
			allReleases = append(allReleases, newSingles...)
			u.mu.Unlock()
		}()

		// waiting and creating all releases in database
		wg.Wait()
		err := u.CreateMultiArtistsAndReleases(allReleases)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("releases are updated")
	}

	// update covers after adding releases
	err := u.UpdateCoversInDB()
	if err != nil {
		log.Fatal(err)
	}
}
