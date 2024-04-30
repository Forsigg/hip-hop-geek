package releases

import (
	"log"

	"hip-hop-geek/internal/models"
	"hip-hop-geek/pkg/covers"
)

func ConvertPostsToReleases(
	posts []models.Post,
	withCover bool,
	coverBook CoverBook,
	releaseType models.ReleaseType,
) []models.Release {
	releases := make([]models.Release, 0)

	for _, post := range posts {

		coverUrl := models.CoverUrl{}
		if withCover {
			log.Println("getting cover for post")
			coverCh := make(chan *covers.Cover)
			go func() {
				cover := coverBook.GetCoverByQuery(post.Query(), 500)
				coverCh <- cover
			}()
			cover := <-coverCh

			if cover.Valid {
				coverUrl.Value = cover.Url
				coverUrl.IsValid = true
			}

		}
		releases = append(releases, models.Release{
			Id: post.Id,
			Artist: models.Artist{
				Name: post.Artist(),
			},
			Type:     releaseType,
			Title:    post.Title(),
			OutDate:  post.ReleaseDate(),
			CoverUrl: coverUrl,
		})
	}

	return releases
}
