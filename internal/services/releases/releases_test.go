package releases

import (
	"reflect"
	"testing"
	"time"

	"hip-hop-geek/internal/models"
	"hip-hop-geek/internal/types"
	"hip-hop-geek/pkg/covers"
)

type StubCoverBook struct{}

func (s *StubCoverBook) GetCoverByQuery(query string, size int) *covers.Cover {
	return &covers.Cover{
		Url:   "cover",
		Valid: true,
	}
}

func TestConvertPostsToReleases(t *testing.T) {
	coverBook := StubCoverBook{}
	cases := []struct {
		name      string
		withCover bool
		posts     []models.Post
		releases  []models.Release
	}{
		{
			"convert without cover",
			false,
			[]models.Post{
				models.NewPost(
					1,
					"21 Savage - American Dream",
					types.NewCustomDate(2024, time.January, 12),
				),
			},
			[]models.Release{
				{
					1,
					models.Artist{"21 Savage"},
					"American Dream",
					models.Album,
					types.NewCustomDate(2024, time.January, 12),
					models.CoverUrl{},
				},
			},
		},
		{
			"convert with cover",
			true,
			[]models.Post{
				models.NewPost(
					1,
					"21 Savage - American Dream",
					types.NewCustomDate(2024, time.January, 12),
				),
			},
			[]models.Release{
				{
					1,
					models.Artist{"21 Savage"},
					"American Dream",
					models.Album,
					types.NewCustomDate(2024, time.January, 12),
					models.CoverUrl{"cover", true},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertPostsToReleases(tc.posts, tc.withCover, &coverBook, models.Album)

			if len(result) != len(tc.releases) {
				t.Fatalf("want %d len releases, got %d", len(tc.releases), len(result))
			}

			if !reflect.DeepEqual(result, tc.releases) {
				t.Errorf("not valid convert posts: want %v got %v", tc.releases, result)
			}
		})
	}
}
