package models

type MessageIdType int

const (
	ReleasesMessage = iota
	TodayReleasesMessage
)

type User struct {
	Id                     int64
	Username               string
	IsTodaySubscribe       bool
	ReleasesMessageId      int64
	ReleasesPageCount      int
	TodayReleasesMessageId int64
	TodayReleasesPageCount int
}
