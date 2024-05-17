package releases

import (
	"reflect"
	"testing"
	"time"

	"hip-hop-geek/internal/db"
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
	cases := []struct {
		name     string
		posts    []models.Post
		releases []models.Release
	}{
		{
			"convert without cover",
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertPostsToReleases(tc.posts, models.Album)

			if len(result) != len(tc.releases) {
				t.Fatalf("want %d len releases, got %d", len(tc.releases), len(result))
			}

			if !reflect.DeepEqual(result, tc.releases) {
				t.Errorf("not valid convert posts: want %v got %v", tc.releases, result)
			}
		})
	}
}

func TestConvertDbReleaseToModelRelease(t *testing.T) {
	cases := []struct {
		name       string
		dbReleases []*db.ReleaseDB
		releases   []models.Release
	}{
		{
			"success case",
			[]*db.ReleaseDB{
				{
					Id:       1,
					Artist:   db.ArtistDB{Id: 1, Name: "21 Savage"},
					Title:    "American Dream",
					Type:     models.Album,
					OutYear:  2024,
					OutMonth: 1,
					OutDay:   12,
					CoverUrl: "",
				},
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
			"with cover",
			[]*db.ReleaseDB{
				{
					Id:       1,
					Artist:   db.ArtistDB{Id: 1, Name: "21 Savage"},
					Title:    "American Dream",
					Type:     models.Album,
					OutYear:  2024,
					OutMonth: 1,
					OutDay:   12,
					CoverUrl: "https://cover.com/album/123",
				},
			},

			[]models.Release{
				{
					1,
					models.Artist{"21 Savage"},
					"American Dream",
					models.Album,
					types.NewCustomDate(2024, time.January, 12),
					models.CoverUrl{
						IsValid: true,
						Value:   "https://cover.com/album/123",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertDbReleaseToModelRelease(tc.dbReleases)

			if len(result) != len(tc.releases) {
				t.Fatalf("want %d len releases, got %d", len(tc.releases), len(result))
			}

			if !reflect.DeepEqual(result, tc.releases) {
				t.Errorf("not valid convert db releases: want %v got %v", tc.releases, result)
			}
		})
	}
}
