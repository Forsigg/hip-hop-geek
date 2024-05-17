package covers

import (
	"fmt"
	"strings"
)

type CoverBook struct {
	fetcher *Fetcher
}

func NewCoverBook() *CoverBook {
	return &CoverBook{
		fetcher: NewFetcher(),
	}
}

type Cover struct {
	Url   string
	Valid bool
}

func (c *CoverBook) GetCoverByQuery(query string, size int) *Cover {
	var cover Cover
	coverResp := c.fetcher.getCover(query)

	if len(coverResp.Results) == 0 {
		return &cover
	}

	if coverResp.Results[0].ArtworkUrl == "" {
		return &cover
	}

	sizedUrl := strings.TrimSuffix(coverResp.Results[0].ArtworkUrl, "100x100bb.jpg")
	sizedUrl += fmt.Sprintf("%dx%d.jpg", size, size)
	return &Cover{
		Url:   sizedUrl,
		Valid: true,
	}
}
