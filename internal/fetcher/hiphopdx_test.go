package fetcher

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"hip-hop-geek/internal/models"
	"hip-hop-geek/internal/types"
)

var isRequestDo = false

const emptyTestJson = `
    {
		"data": {
		    "posts": []
		}
	}`

const testJson = `{
                "data": {
                    "posts": [
                        {
                            "ID": 1,
                            "post_title": "21 Savage - American Dream",
                            "post_date": "2024-01-12 00:00:00"
                        }
                    ]
                }
            }`

type StubHttpClient struct {
	respBodyFull  string
	respBodyEmpty string
}

func (s *StubHttpClient) Do(req *http.Request) (*http.Response, error) {
	log.Println("do mock request...")
	if isRequestDo {
		resp := httptest.NewRecorder()
		resp.Body.Write([]byte(s.respBodyEmpty))
		return resp.Result(), nil
	} else {

		resp := httptest.NewRecorder()
		resp.Body.Write([]byte(s.respBodyFull))
		isRequestDo = true

		return resp.Result(), nil
	}
}

func TestGetSinglesPosts(t *testing.T) {
	fetcher := HipHopDXFetcher{
		Client: &StubHttpClient{
			respBodyEmpty: emptyTestJson,
			respBodyFull:  testJson,
		},
	}

	isRequestDo = false
	t.Run("happy path", func(t *testing.T) {
		expected := []models.Post{
			models.NewPost(
				1,
				"21 Savage - American Dream",
				types.NewCustomDate(2024, time.January, 12),
			),
		}

		got, err := fetcher.GetSinglesPosts(2024)
		assert.NoError(t, err)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("incorrect posts: want %v, got %v", expected, got)
		}
	})
}

func TestReleasesPosts(t *testing.T) {
	fetcher := HipHopDXFetcher{
		Client: &StubHttpClient{
			respBodyEmpty: emptyTestJson,
			respBodyFull:  testJson,
		},
	}

	isRequestDo = false
	t.Run("happy path", func(t *testing.T) {
		expected := []models.Post{
			models.NewPost(
				1,
				"21 Savage - American Dream",
				types.NewCustomDate(2024, time.January, 12),
			),
		}

		got, err := fetcher.GetReleasesPosts(2024, time.January)
		assert.NoError(t, err)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("incorrect posts: want %v, got %v", expected, got)
		}
	})
}

func TestBuildReleaseUrl(t *testing.T) {
	fetcher := HipHopDXFetcher{
		Client: &StubHttpClient{
			respBodyEmpty: emptyTestJson,
			respBodyFull:  testJson,
		},
		currentReq: nil,
	}
	cases := []struct {
		name     string
		year     int
		month    time.Month
		expected string
	}{
		{
			"2024 January test",
			2024,
			time.January,
			`https://app.hiphopdx.com/wp-json/hiphopdx-api/v1/get_posts?monthnum=1&paged=1&post_status=publish%2Cfuture&post_type=release-date&posts_per_page=99&year=2024`,
		},
		{
			"2023 March test",
			2023,
			time.March,
			`https://app.hiphopdx.com/wp-json/hiphopdx-api/v1/get_posts?monthnum=3&paged=1&post_status=publish%2Cfuture&post_type=release-date&posts_per_page=99&year=2023`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fetcher.buildReleasesUrl(tc.year, tc.month, 1)
			if got != tc.expected {
				t.Errorf("url not correct build: want '%s' got '%s'", tc.expected, got)
			}
		})
	}
}

func TestBuildSinglesUrl(t *testing.T) {
	fetcher := HipHopDXFetcher{
		Client: &StubHttpClient{
			respBodyEmpty: emptyTestJson,
			respBodyFull:  testJson,
		},
		currentReq: nil,
	}
	cases := []struct {
		name     string
		year     int
		page     int
		expected string
	}{
		{
			"2024 first page",
			2024,
			1,
			"https://app.hiphopdx.com/wp-json/hiphopdx-api/v1/get_posts?paged=1&post_status=publish%2Cfuture&post_type=single&posts_per_page=100&year=2024",
		},
		{
			"2023 third page",
			2023,
			3,
			"https://app.hiphopdx.com/wp-json/hiphopdx-api/v1/get_posts?paged=3&post_status=publish%2Cfuture&post_type=single&posts_per_page=100&year=2023",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fetcher.buildSinglesUrl(tc.year, 0, tc.page)

			if got != tc.expected {
				t.Errorf("not equal urls: want '%s' got '%s'", tc.expected, got)
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	fetcher := HipHopDXFetcher{
		Client: &StubHttpClient{
			respBodyEmpty: emptyTestJson,
			respBodyFull:  testJson,
		},
		currentReq: nil,
	}
	cases := []struct {
		name     string
		json     string
		expected []models.Post
	}{
		{
			"happy path",
			testJson,
			[]models.Post{
				models.NewPost(
					1,
					"21 Savage - American Dream",
					types.NewCustomDate(2024, time.January, 12),
				),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := getJsonResponse(t, c.json)
			result := fetcher.parseResponse(resp)

			if !reflect.DeepEqual(result, c.expected) {
				t.Errorf("expected %v, got %v", c.expected, result)
			}
		})
	}
}

func getJsonResponse(t testing.TB, json string) *http.Response {
	t.Helper()
	resp := httptest.NewRecorder()
	resp.Body.Write([]byte(json))
	return resp.Result()
}
