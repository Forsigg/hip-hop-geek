package fetcher

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"hip-hop-geek/internal/models"
)

const (
	todayHipHopHistoryUrl = "https://todayinhiphophistory.com"

	aClassMediaPhotoImageSelector = ".post_media_photo_anchor"
	divClassText                  = "div.caption"
	divPostSelector               = "div.post"
	dateLinkSelector              = "div.date>a"

	dateLayout = "Jan. 02 2006"

	twoDaysMinutes = float64((60 * 24) * 2)
)

var (
	ErrImageUrlNotFound = errors.New("image url not found")
	ErrTextPostNotFound = errors.New("text post not found")
	ErrPostsNotFound    = errors.New("posts not found")
)

type TodayHipHopFetcher struct {
	Client     CustomHttpClient
	currentReq *http.Request
}

func NewTodayHipHopFetcher() *TodayHipHopFetcher {
	return &TodayHipHopFetcher{
		&http.Client{},
		nil,
	}
}

func (f *TodayHipHopFetcher) Close() {
	if f.currentReq != nil {
		f.currentReq.Body.Close()
	}
}

func (f *TodayHipHopFetcher) GetTodayEvents() ([]*models.TodayPost, error) {
	htmlBody := f.getHTML()
	doc := f.parseResponse(htmlBody)
	post, err := f.getPostsFromDoc(doc, time.Now().UTC())

	return post, err
}

func (f *TodayHipHopFetcher) getHTML() io.ReadCloser {
	req, err := http.NewRequest(http.MethodGet, todayHipHopHistoryUrl, nil)
	if err != nil {
		log.Fatalf("error while creating request for today history site: %s", err)
	}
	f.currentReq = req
	resp, err := f.Client.Do(req)
	f.currentReq = nil
	if err != nil {
		log.Fatalf("error while getting today history html: %s", err)
	}

	return resp.Body
}

func (f *TodayHipHopFetcher) parseResponse(htmlBody io.ReadCloser) *goquery.Document {
	defer htmlBody.Close()
	doc, err := goquery.NewDocumentFromReader(htmlBody)
	if err != nil {
		log.Fatalf("error while parsing html: %s", err)
	}

	return doc
}

func (f *TodayHipHopFetcher) getPostsFromDoc(
	doc *goquery.Document,
	now time.Time,
) ([]*models.TodayPost, error) {
	if time.Now().UTC().Sub(now).Minutes() > twoDaysMinutes {
		return nil, ErrPostsNotFound
	}

	var posts []*models.TodayPost
	doc.Find(divPostSelector).Each(func(i int, s *goquery.Selection) {
		date := s.Find(dateLinkSelector).Text()
		tt, err := time.Parse(dateLayout, date)
		if err != nil {
			log.Fatal(err)
		}

		if tt.Month() == now.Month() && tt.Day() == now.Day() && tt.Year() == now.Year() {
			text := s.Find(divClassText).Find("p").Text()
			text = strings.TrimPrefix(text, "Today in Hip Hop History:")
			image, _ := s.Find(aClassMediaPhotoImageSelector).Attr("data-big-photo")

			posts = append(posts, &models.TodayPost{
				Text: text,
				Url:  image,
			})
		}
	})

	if len(posts) == 0 {
		return f.getPostsFromDoc(doc, now.AddDate(0, 0, -1))
	}

	return posts, nil
}
