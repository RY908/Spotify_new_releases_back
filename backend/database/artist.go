package database

import (
	_ "github.com/go-sql-driver/mysql"
)

type ArtistInfo struct {
    ArtistId	string 		`db:"artistId" 	json:"artistId"`
	Name 		string 		`db:"name" 		json:"name"`
	Url 		string 		`db:"url" 		json:"url"`
	IconUrl		string 		`db:iconUrl 	json:"iconUrl"`
}

type ArtistRes struct {
	ArtistId	string 		`json:"artistId"`
	Name 		string 		`json:"name"`
	Url 		string 		`json:"url"`
	IconUrl		string 		`json:"iconUrl"`
	IfFollowing bool		`json:"ifFollowing"`
}

// insert artist in database
func (d *MyDbMap) InsertArtist(artistId, name, url, iconUrl string) error {
	count, err := d.DbMap.SelectInt("select count(*) from Artist where artistId = ?", artistId) // check if artist already exists in database
	if err != nil {
		return err
	}

	// if artist does not exist, then insert artist into database
	if count == 0 {
		err = d.DbMap.Insert(&ArtistInfo{ArtistId:artistId, Name:name, Url:url, IconUrl:iconUrl})
	}
	
	if err != nil {
		return err
	}
	return nil
}

// insert multiple artists in database
func (d *MyDbMap) InsertArtists(artists []ArtistInfo) error {
	trans, err := d.DbMap.Begin()
	if err != nil {
		return err	
	}
	
	for _, artist := range artists {
		artistId := artist.ArtistId
		count, err := trans.SelectInt("select count(*) from Artist where artistId = ?", artistId) // check if artist already exists in database
		if err != nil {
			return err
		}

		// if artist does not exist, then insert artist into database
		if count == 0 {
			err = trans.Insert(&artist)
		}
		
		if err != nil {
			return err
		}
	}

	return trans.Commit()
}

// get artists that the user listened to or follows
func (d *MyDbMap) GetArtistsFromUserId(userId string) ([]ArtistRes, error) {
	var artists []ArtistRes
	cmd := "select Artist.artistId, Artist.name, Artist.url, Artist.iconUrl, ListenTo.ifFollowing from Artist inner join ListenTo on Artist.artistId = ListenTo.artistId where ListenTo.userId = ? and ListenTo.listenCount >= 2"
	_, err := d.DbMap.Select(&artists, cmd, userId)
	if err != nil {
		return nil, err  
	}
	return artists, nil
}
