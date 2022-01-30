package schema

type Artist struct {
	ID      string `db:"artistId"`
	Name    string `db:"name"`
	Url     string `db:"url"`
	IconUrl string `db:"iconUrl"`
}
