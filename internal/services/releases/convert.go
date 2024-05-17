package releases

import (
	"log"
	"strings"
	"time"

	"hip-hop-geek/internal/db"
	"hip-hop-geek/internal/models"
	"hip-hop-geek/internal/types"
)

func ConvertPostsToReleases(
	posts []models.Post,
	releaseType models.ReleaseType,
) []models.Release {
	releases := make([]models.Release, 0, len(posts))

	for _, post := range posts {
		// skip releases without divider symbol
		if !strings.Contains(post.QueryField, " - ") && !strings.Contains(post.QueryField, " – ") {
			log.Printf("release %s skipped", post.QueryField)
			continue
		}

		releases = append(releases, models.Release{
			Id: post.Id,
			Artist: models.Artist{
				Name: post.Artist(),
			},
			Type:     releaseType,
			Title:    post.Title(),
			OutDate:  post.ReleaseDate(),
			CoverUrl: models.CoverUrl{},
		})
	}

	return releases
}

func ConvertDbReleaseToModelRelease(dbReleases []*db.ReleaseDB) []models.Release {
	releases := make([]models.Release, 0, len(dbReleases))

	for _, dbRelease := range dbReleases {

		coverUrl := models.CoverUrl{}
		if dbRelease.CoverUrl != "" {
			coverUrl.Value = dbRelease.CoverUrl
			coverUrl.IsValid = true
		}

		releases = append(releases, models.Release{
			Id: dbRelease.Id,
			Artist: models.Artist{
				Name: dbRelease.Artist.Name,
			},
			Type:  models.ReleaseType(dbRelease.Type),
			Title: dbRelease.Title,
			OutDate: types.NewCustomDate(
				dbRelease.OutYear, time.Month(dbRelease.OutMonth), dbRelease.OutDay,
			),
			CoverUrl: coverUrl,
		})
	}

	return releases
}
