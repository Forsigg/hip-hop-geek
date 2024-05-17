package fetcher

import (
	"io"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"hip-hop-geek/internal/models"
)

const htmlBody = `
        <html>
            <body>
                <h1>hello world!</h1>
                <a class="post_media_photo_anchor" data-big-photo="https://youfool.com"></a>
                <div class="caption">
                    <p>Today in Hip Hop History:</p>
                    <p>Hip Hop was born 11 August 1973</p>
                </div>
            </body>
        </html>`

func TestGetHTML(t *testing.T) {
	isRequestDo = false

	t.Run("happy path", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  htmlBody,
				respBodyEmpty: "",
			},
			nil,
		}
		want := []byte(htmlBody)
		got, err := io.ReadAll(fetcher.getHTML())
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	isRequestDo = false
	t.Run("empty response", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  "",
				respBodyEmpty: "",
			},
			nil,
		}
		want := []byte("")
		got, err := io.ReadAll(fetcher.getHTML())
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestParseEventsResponse(t *testing.T) {
	isRequestDo = false
	t.Run("success case", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  htmlBody,
				respBodyEmpty: "",
			},
			nil,
		}
		htmlReader := strings.NewReader(htmlBody)
		want, _ := goquery.NewDocumentFromReader(htmlReader)

		htmlB := fetcher.getHTML()
		got := fetcher.parseResponse(htmlB)

		assert.Equal(t, want, got)
	})

	isRequestDo = false
	t.Run("empty body", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  "",
				respBodyEmpty: "",
			},
			nil,
		}
		htmlReader := strings.NewReader("")
		want, _ := goquery.NewDocumentFromReader(htmlReader)

		htmlB := fetcher.getHTML()
		got := fetcher.parseResponse(htmlB)

		assert.Equal(t, want, got)
	})
}

func TestGetImageUrl(t *testing.T) {
	isRequestDo = false
	t.Run("success case", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  htmlBody,
				respBodyEmpty: "",
			},
			nil,
		}
		want := "https://youfool.com"

		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getImageUrl(doc)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	isRequestDo = false
	t.Run("error if image not found", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  "<html><body><h1>Hello world!</h1></body></html>",
				respBodyEmpty: "",
			},
			nil,
		}

		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getImageUrl(doc)

		assert.ErrorIs(t, err, ErrImageUrlNotFound)
		assert.Equal(t, "", got)
	})
}

func TestGetTextPost(t *testing.T) {
	isRequestDo = false
	t.Run("success case", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  htmlBody,
				respBodyEmpty: "",
			},
			nil,
		}
		want := "Hip Hop was born 11 August 1973"

		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getEventText(doc)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	isRequestDo = false
	t.Run("error if text post not found", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  "<html><body><h1>Hello world!</h1></body></html>",
				respBodyEmpty: "",
			},
			nil,
		}

		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getEventText(doc)

		assert.ErrorIs(t, err, ErrTextPostNotFound)
		assert.Equal(t, "", got)
	})
}

func TestGetPostFromDoc(t *testing.T) {
	isRequestDo = false

	t.Run("success case", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  htmlBody,
				respBodyEmpty: "",
			},
			nil,
		}
		want := &models.TodayPost{
			Text: "Hip Hop was born 11 August 1973",
			Url:  "https://youfool.com",
		}

		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getPostFromDoc(doc)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	isRequestDo = false
	t.Run("image url not found", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull:  "<html><body><h1>Hello world!</h1></body></html>",
				respBodyEmpty: "",
			},
			nil,
		}
		var want *models.TodayPost
		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getPostFromDoc(doc)

		assert.ErrorIs(t, err, ErrImageUrlNotFound)
		assert.Equal(t, want, got)
	})

	isRequestDo = false
	t.Run("image url not found", func(t *testing.T) {
		fetcher := TodayHipHopFetcher{
			&StubHttpClient{
				respBodyFull: `<html>
                                    <body>
                                        <h1>Hello world!</h1>
                                        <a class="post_media_photo_anchor" data-big-photo="https://youfool.com"></a>
                                    </body>
                                </html>`,
				respBodyEmpty: "",
			},
			nil,
		}
		var want *models.TodayPost
		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getPostFromDoc(doc)

		assert.ErrorIs(t, err, ErrTextPostNotFound)
		assert.Equal(t, want, got)
	})
}
