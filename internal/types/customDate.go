package types

import (
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

func NewCustomDate(year int, month time.Month, day int) CustomDate {
	return CustomDate{
		time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	layout := "2006-01-02 15:04:05"

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}

	c.Time, err = time.Parse(layout, s)
	return
}
