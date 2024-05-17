package sqlite

import "errors"

var (
	ErrArtistAlreadyExists  = errors.New("artist with this name already exists")
	ErrReleaseAlreadyExists = errors.New("release with that id already exists")
	ErrReleasesNotFound     = errors.New("releases not found")
)
