package fetcher

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"hip-hop-geek/internal/models"
)

const (
	todayHipHopHistoryUrl = "https://todayinhiphophistory.com"

	aClassMediaPhotoImageSelector = ".post_media_photo_anchor"
	divClassText                  = "div.caption"
)

var (
	ErrImageUrlNotFound = errors.New("image url not found")
	ErrTextPostNotFound = errors.New("text post not found")
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

func (f *TodayHipHopFetcher) GetTodayEvent() (*models.TodayPost, error) {
	htmlBody := f.getHTML()
	doc := f.parseResponse(htmlBody)
	post, err := f.getPostFromDoc(doc)

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

func (f *TodayHipHopFetcher) getImageUrl(doc *goquery.Document) (string, error) {
	node := doc.Find(aClassMediaPhotoImageSelector).First()
	val, exists := node.Attr("data-big-photo")
	if exists {
		return val, nil
	}

	return "", ErrImageUrlNotFound
}

func (f *TodayHipHopFetcher) getEventText(doc *goquery.Document) (string, error) {
	node := doc.Find(divClassText).First()
	text := node.Find("p").Text()
	if text == "" {
		return "", ErrTextPostNotFound
	}

	text = strings.TrimPrefix(text, "Today in Hip Hop History:")
	return text, nil
}

func (f *TodayHipHopFetcher) getPostFromDoc(doc *goquery.Document) (*models.TodayPost, error) {
	imageUrl, err := f.getImageUrl(doc)
	if err != nil {
		return nil, fmt.Errorf("error while getting image url: %w", err)
	}

	postText, err := f.getEventText(doc)
	if err != nil {
		return nil, fmt.Errorf("error while getting post text: %w", err)
	}

	return &models.TodayPost{
		Text: postText,
		Url:  imageUrl,
	}, nil
}
