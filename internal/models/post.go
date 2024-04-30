package models

import (
	"strings"

	"hip-hop-geek/internal/types"
)

type Post struct {
	Id               int              `json:"ID"`
	QueryField       string           `json:"post_title"`
	ReleaseDateField types.CustomDate `json:"post_date"`
}

func NewPost(id int, query string, releaseDate types.CustomDate) Post {
	return Post{
		Id:               id,
		QueryField:       query,
		ReleaseDateField: releaseDate,
	}
}

func (p Post) Artist() string {
	divider := "-"
	if strings.Contains(p.QueryField, "–") {
		divider = "–"
	}

	return strings.TrimSuffix(strings.Split(p.QueryField, divider)[0], " ")
}

func (p Post) Title() string {
	divider := "-"
	if strings.Contains(p.QueryField, "–") {
		divider = "–"
	}
	return strings.TrimPrefix(strings.Split(p.QueryField, divider)[1], " ")
}

func (p Post) Query() string {
	return p.QueryField
}

func (p Post) ReleaseDate() types.CustomDate {
	return p.ReleaseDateField
}
