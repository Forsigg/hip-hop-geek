package db

import "hip-hop-geek/internal/models"

type ReleaseRepositoryInterface interface {
	AddRelease(release models.Release, artId int) (int, error)
	GetReleaseById(id int) (*ReleaseDB, error)
	GetReleaseByTitle(title string) (*ReleaseDB, error)
}

type ArtistsRepositoryInterface interface {
	AddArtist(artistName string) (int, error)
	GetArtistByName(artistName string) (*ArtistDB, error)
	GetArtistById(id int) (*ArtistDB, error)
}

type DbRepository struct {
	ReleaseRepositoryInterface
	ArtistsRepositoryInterface
}
