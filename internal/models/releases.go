package models

import (
	"fmt"

	"hip-hop-geek/internal/types"
)

type ReleaseType int

const (
	Album = iota + 1
	Single
)

type CoverUrl struct {
	Value   string
	IsValid bool
}

type Artist struct {
	Name string
}

type Release struct {
	Id       int
	Artist   Artist
	Title    string
	Type     ReleaseType
	OutDate  types.CustomDate
	CoverUrl CoverUrl
}

func (r Release) String() string {
	return fmt.Sprintf("%s - %s (out %s)", r.Artist.Name, r.Title, r.OutDate)
}
