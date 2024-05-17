package db

type ArtistDB struct {
	Id   int
	Name string
}

type ReleaseDB struct {
	Id       int
	Artist   ArtistDB
	Title    string
	Type     int
	OutYear  int
	OutMonth int
	OutDay   int
	CoverUrl string
}
