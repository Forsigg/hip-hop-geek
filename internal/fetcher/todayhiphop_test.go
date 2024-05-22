package fetcher

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"hip-hop-geek/internal/models"
)

const htmlBody = `
        <html>
            <body>
                <h1>hello world!</h1>
                <div class="post text">
                    <a class="post_media_photo_anchor" data-big-photo="https://youfool.com"></a>
                    <div class="caption">
                        <p>Today in Hip Hop History:</p>
                        <p>Hip Hop was born 11 August 1973</p>
                    </div>
                    <div class="date">
                        <a href="123123">May. 21 2024</a>
                    </div>
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
		want := []*models.TodayPost{
			{
				Text: "Hip Hop was born 11 August 1973",
				Url:  "https://youfool.com",
			},
		}

		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getPostFromDoc(doc, time.Now().UTC())

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
		var want []*models.TodayPost
		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getPostFromDoc(doc, time.Now().UTC())

		assert.ErrorIs(t, err, ErrPostsNotFound)
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
		var want []*models.TodayPost
		htmlB := fetcher.getHTML()
		doc := fetcher.parseResponse(htmlB)
		got, err := fetcher.getPostFromDoc(doc, time.Now().UTC())

		assert.ErrorIs(t, err, ErrPostsNotFound)
		assert.Equal(t, want, got)
	})
}
