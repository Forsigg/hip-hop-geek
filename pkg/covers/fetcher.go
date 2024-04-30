package covers

import (
	"encoding/json"
	"net/http"
)

type Fetcher struct {
	base_url string
}

type result struct {
	CollectionType string `json:"collectionType"`
	ArtistId       int    `json:"artistId"`
	ArtistName     string `json:"artistName"`
	CollectionId   int    `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	ArtworkUrl     string `json:"artworkUrl100"`
}

type coverResponse struct {
	ResultCount int      `json:"resultCount"`
	Results     []result `json:"results"`
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		base_url: "https://itunes.apple.com/search",
	}
}

func (f *Fetcher) getCover(query string) *coverResponse {
	resp := f.doRequest(query)

	var respModel coverResponse
	json.NewDecoder(resp.Body).Decode(&respModel)

	return &respModel
}

func (f *Fetcher) generateRequest(query string) *http.Request {
	payload := map[string]string{
		"term":   query,
		"limit":  "1",
		"entity": "musicArtist,musicTrack,album,mix,song",
		"media":  "music",
	}
	req, _ := http.NewRequest(http.MethodGet, f.base_url, nil)
	q := req.URL.Query()
	for param, value := range payload {
		q.Add(param, value)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func (f *Fetcher) doRequest(query string) *http.Response {
	req := f.generateRequest(query)

	client := http.Client{}
	resp, _ := client.Do(req)
	return resp
}
