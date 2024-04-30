package fetcher

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"hip-hop-geek/internal/models"
)

const (
	postsUrl = "https://app.hiphopdx.com/wp-json/hiphopdx-api/v1/get_posts"
)

type CustomHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HipHopDXFetcher struct {
	Client CustomHttpClient
}

type PostsData struct {
	Posts []models.Post `json:"posts"`
}

type ReleaseResponse struct {
	Data PostsData `json:"data"`
}

func NewHipHopDXFetcher() *HipHopDXFetcher {
	return &HipHopDXFetcher{
		&http.Client{},
	}
}

func (h *HipHopDXFetcher) GetSinglesPosts(year int) []models.Post {
	log.Println("start getting singles posts...")
	posts := make([]models.Post, 0)
	for page := 1; true; page++ {
		url := h.buildSinglesUrl(year, 0, 1)
		resp := h.DoRequest(url)

		respPosts := h.parseResponse(resp)
		if len(respPosts) == 0 {
			break
		}

		posts = append(posts, respPosts...)
	}

	return posts
}

func (h *HipHopDXFetcher) GetReleasesPosts(
	year int,
	month time.Month,
) []models.Post {
	posts := make([]models.Post, 0)

	for page := 1; true; page++ {
		url := h.buildReleasesUrl(year, month, page)
		resp := h.DoRequest(url)
		postsBatch := h.parseResponse(resp)

		if len(postsBatch) == 0 {
			break
		}
		posts = append(posts, postsBatch...)

	}

	return posts
}

func (h *HipHopDXFetcher) parseResponse(resp *http.Response) []models.Post {
	log.Println("parsing response...")
	var p ReleaseResponse
	err := json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		log.Fatalf("error while parsing response: %s", err)
	}

	return p.Data.Posts
}

func (h *HipHopDXFetcher) buildReleasesUrl(year int, month time.Month, page int) string {
	queries := map[string][]string{
		"post_type":      {"release-date"},
		"posts_per_page": {"99"},
		"paged":          {strconv.Itoa(page)},
		"post_status":    {"publish,future"},
		"year":           {strconv.Itoa(year)},
	}
	if month != 0 {
		queries["monthnum"] = []string{strconv.Itoa(int(month))}
	}

	return BuildUrl(queries)
}

func (h *HipHopDXFetcher) buildSinglesUrl(year int, month time.Month, page int) string {
	queries := map[string][]string{
		"post_type":      {"single"},
		"posts_per_page": {"100"},
		"paged":          {strconv.Itoa(page)},
		"post_status":    {"publish,future"},
		"year":           {strconv.Itoa(year)},
	}
	return BuildUrl(queries)
}

func (h *HipHopDXFetcher) DoRequest(url string) *http.Response {
	log.Printf("do request to %s", url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := h.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
