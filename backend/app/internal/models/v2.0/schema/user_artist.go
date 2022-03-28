package schema

type UserArtist struct {
	ID          string `db:"artistId"`
	Name        string `db:"name"`
	Url         string `db:"url"`
	IconUrl     string `db:"iconUrl"`
	IsFollowing bool   `db:"ifFollowing"`
}
