package db

import (
	"time"

	"hip-hop-geek/internal/models"
)

type ReleaseRepositoryInterface interface {
	AddRelease(release models.Release, artId int) (int, error)
	GetReleaseById(id int) (*ReleaseDB, error)
	GetReleaseByTitle(title string) (*ReleaseDB, error)
	GetReleasesByMonth(month time.Month, year, limit, offset int) ([]*ReleaseDB, error)
	GetReleasesByYear(year, limit, offset int) ([]*ReleaseDB, error)
	GetReleasesByDay(year int, month time.Month, day, limit, offset int) ([]*ReleaseDB, error)
	GetReleasesWithoutCover() ([]*ReleaseDB, error)
	UpdateReleaseCoverUrl(releaseId int, coverUrl string) error
	CloseReleaseRepo()
}

type ArtistsRepositoryInterface interface {
	AddArtist(artistName string) (int, error)
	GetArtistByName(artistName string) (*ArtistDB, error)
	GetArtistById(id int) (*ArtistDB, error)
	CloseArtistRepo()
}

type UsersRepositoryInterface interface {
	AddUser(user models.User) error
	GetAllSubscribers() ([]*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	SetTodaySubscribe(userId int64, isSubscribe bool) error
	SetUserState(userId int64, messageType, messageId int, pageCount int) error
}

type DbRepository interface {
	ReleaseRepositoryInterface
	ArtistsRepositoryInterface
	UsersRepositoryInterface
	Close()
}
