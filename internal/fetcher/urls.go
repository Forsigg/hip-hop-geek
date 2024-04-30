package fetcher

import (
	"log"
	"net/url"
)

func BuildUrl(queries map[string][]string) string {
	url, err := url.Parse(postsUrl)
	if err != nil {
		log.Fatal(err)
	}
	AddQueriesToUrl(url, queries)

	return url.String()
}

func AddQueriesToUrl(url_link *url.URL, queries map[string][]string) {
	url_link.RawQuery = url.Values(queries).Encode()
}
