package releases

import (
	"errors"
	"log"
	"time"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/db/sqlite"
	"hip-hop-geek/internal/models"
	"hip-hop-geek/pkg/covers"
)

const AllMonths = 0

type CoverBook interface {
	GetCoverByQuery(query string, size int) *covers.Cover
}

type EventsFetcher interface {
	GetTodayEvents() ([]*models.TodayPost, error)
	Close()
}

type ReleaseFetcher interface {
	GetSinglesPosts(year int) ([]models.Post, error)
	GetReleasesPosts(year int, month time.Month) ([]models.Post, error)
	Close()
}

type HipHopService struct {
	db.DbRepository
	ReleaseFetcher ReleaseFetcher
	EventsFetcher  EventsFetcher
}

func NewHipHopService(
	dbRepo db.DbRepository,
	releasesFetcher ReleaseFetcher,
	eventsFetcher EventsFetcher,
) *HipHopService {
	return &HipHopService{
		dbRepo,
		releasesFetcher,
		eventsFetcher,
	}
}

func (h *HipHopService) Close() {
	h.DbRepository.Close()
	h.ReleaseFetcher.Close()
	h.EventsFetcher.Close()
}

// func (h *HipHopService) AddUser(user models.User) error {
// 	return h.Repo.AddUser(user)
// }

// func (h *HipHopService) SetTodaySubscribe(
// 	user models.User,
// 	isSubscribe bool,
// ) (*models.User, error) {
// 	return h.Repo.SetTodaySubscribe(user, isSubscribe)
// }

func (h *HipHopService) GetTodayEvents() ([]*models.TodayPost, error) {
	return h.EventsFetcher.GetTodayEvents()
}

func (h *HipHopService) FetchReleases(year int) ([]models.Release, error) {
	releases := make([]models.Release, 0, 5)
	for monthNum := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12} {
		posts, err := h.ReleaseFetcher.GetReleasesPosts(year, time.Month(monthNum))
		if err != nil {
			return nil, err
		}
		convertedPosts := ConvertPostsToReleases(posts, models.Album)
		releases = append(releases, convertedPosts...)
	}

	return releases, nil
}

func (h *HipHopService) FetchSingles(year int) ([]models.Release, error) {
	singlesPosts, err := h.ReleaseFetcher.GetSinglesPosts(year)
	if err != nil {
		return nil, err
	}
	return ConvertPostsToReleases(singlesPosts, models.Single), nil
}

func (h *HipHopService) GetMonthReleases(
	year int,
	month time.Month,
	limit,
	offset int,
) []models.Release {
	releases, err := h.GetReleasesByMonth(month, year, limit, offset)
	if err != nil {
		if errors.Is(err, sqlite.ErrReleasesNotFound) {
			return nil
		} else {
			log.Fatal(err)
		}
	}
	rels := ConvertDbReleaseToModelRelease(releases)
	return rels
}

func (h *HipHopService) GetAllYearReleases(year, limit, offset int) []models.Release {
	releases, err := h.GetReleasesByYear(year, limit, offset)
	if err != nil {
		if errors.Is(err, sqlite.ErrReleasesNotFound) {
			return nil
		} else {
			log.Fatal(err)
		}
	}
	rels := ConvertDbReleaseToModelRelease(releases)
	return rels
}

func (h *HipHopService) GetAllYearSingles(year int, withCover bool) []models.Release {
	rels := make([]models.Release, 0)
	// posts := h.GetSinglesPosts(year)
	// rels := ConvertPostsToReleases(posts, withCover, covers.NewCoverBook(), models.Single)
	return rels
}

func (h *HipHopService) GetReleasesByDay(
	year int,
	month time.Month,
	day, limit, offset int,
) []models.Release {
	releases, err := h.DbRepository.GetReleasesByDay(year, month, day, limit, offset)
	if err != nil {
		if errors.Is(err, sqlite.ErrReleasesNotFound) {
			return nil
		} else {
			log.Fatal(err)
		}
	}

	return ConvertDbReleaseToModelRelease(releases)
}
