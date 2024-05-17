package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
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
	Client     CustomHttpClient
	currentReq *http.Request
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
		nil,
	}
}

func (f *HipHopDXFetcher) Close() {
	if f.currentReq != nil {
		f.currentReq.Body.Close()
	}
}

func (f *HipHopDXFetcher) GetSinglesPosts(year int) ([]models.Post, error) {
	log.Println("start getting singles posts...")
	posts := make([]models.Post, 0)
	for page := 1; true; page++ {
		url := f.buildSinglesUrl(year, 0, page)
		resp, err := f.DoRequest(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("too many requests, get some sleep (1 min)...")
			time.Sleep(60 * time.Second)
			resp, err = f.DoRequest(url)
			if err != nil {
				return nil, err
			}
		}
		respPosts := f.parseResponse(resp)
		if len(respPosts) == 0 {
			break
		}

		posts = append(posts, respPosts...)
	}

	return posts, nil
}

func (f *HipHopDXFetcher) GetReleasesPosts(
	year int,
	month time.Month,
) ([]models.Post, error) {
	posts := make([]models.Post, 0)

	for page := 1; true; page++ {
		url := f.buildReleasesUrl(year, month, page)
		resp, err := f.DoRequest(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("too many requests, get some sleep (1 min)...")
			time.Sleep(60 * time.Second)
			resp, err = f.DoRequest(url)
			if err != nil {
				return nil, err
			}
		}
		postsBatch := f.parseResponse(resp)

		if len(postsBatch) == 0 {
			break
		}
		posts = append(posts, postsBatch...)

	}

	return posts, nil
}

func (f *HipHopDXFetcher) parseResponse(resp *http.Response) []models.Post {
	defer resp.Body.Close()
	log.Println("parsing response...")
	var p ReleaseResponse
	err := json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		bodyStr, _ := io.ReadAll(resp.Body)
		log.Fatalf("error while parsing response: %s \njson body: %s", err, bodyStr)
	}

	return p.Data.Posts
}

func (f *HipHopDXFetcher) buildReleasesUrl(year int, month time.Month, page int) string {
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

func (f *HipHopDXFetcher) buildSinglesUrl(year int, month time.Month, page int) string {
	queries := map[string][]string{
		"post_type":      {"single"},
		"posts_per_page": {"100"},
		"paged":          {strconv.Itoa(page)},
		"post_status":    {"publish,future"},
		"year":           {strconv.Itoa(year)},
	}
	return BuildUrl(queries)
}

func (f *HipHopDXFetcher) DoRequest(url string) (*http.Response, error) {
	log.Printf("do request to %s", url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	f.currentReq = req
	resp, err := f.Client.Do(req)
	f.currentReq = nil
	if err != nil {
		return nil, fmt.Errorf("error while do request: %w", err)
	}

	return resp, nil
}
