package releases

import (
	"time"

	"hip-hop-geek/internal/models"
	"hip-hop-geek/pkg/covers"
)

const AllMonths = 0

type CoverBook interface {
	GetCoverByQuery(query string, size int) *covers.Cover
}

type Fetcher interface {
	GetSinglesPosts(year int) []models.Post
	GetReleasesPosts(year int, month time.Month) []models.Post
}

type HipHopDXService struct {
	Fetcher Fetcher
}

func NewHipHopDXService(fetcher Fetcher) *HipHopDXService {
	return &HipHopDXService{
		fetcher,
	}
}

func (h *HipHopDXService) GetMonthReleases(
	year int,
	month time.Month,
	withCover bool,
) []models.Release {
	posts := h.Fetcher.GetReleasesPosts(year, month)
	rels := ConvertPostsToReleases(posts, withCover, covers.NewCoverBook(), models.Album)
	return rels
}

func (h *HipHopDXService) GetAllYearReleases(year int, withCover bool) []models.Release {
	posts := h.Fetcher.GetReleasesPosts(year, AllMonths)
	rels := ConvertPostsToReleases(posts, withCover, covers.NewCoverBook(), models.Album)
	return rels
}

func (h *HipHopDXService) GetAllYearSingles(year int, withCover bool) []models.Release {
	posts := h.Fetcher.GetSinglesPosts(year)
	rels := ConvertPostsToReleases(posts, withCover, covers.NewCoverBook(), models.Single)
	return rels
}
