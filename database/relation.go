package database

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	"time"
	"fmt"
)

type ListenTo struct {
	ListenId	int64 		`db:listenId`
	ArtistId	string 		`db:artistId`
	UserId		string 		`db:userId`
	ListenCount	int64 		`db:listenCount`
	Timestamp	time.Time	`db:"timestamp"`
	IfFollowing bool 		`db:"ifFollowing"`
}

// insert relation
func (d *MyDbMap) InsertRelation(artistId, userId string, timestamp time.Time, ifFollowing bool) error {
	count, err := d.DbMap.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ?", artistId, userId)
	if err != nil {
		return err
	}
	if count == 0 {
		err = d.DbMap.Insert(&ListenTo{ArtistId:artistId, UserId: userId, Timestamp:timestamp, IfFollowing: ifFollowing})
	}
	if err != nil {
		return err
	}
	return nil
}

// insert multiple relations
func (d *MyDbMap) InsertRelations(artists []ArtistInfo, counter map[string]int, userId string, timestamp time.Time, ifFollowing bool) error {
	trans, err := d.DbMap.Begin()
	if err != nil {
		return err
	}

	for _, artist := range artists {
		artistId := artist.ArtistId
		count, err := trans.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ?", artistId, userId)
		if err != nil {
			return err
		}
		if count == 0 {
			if ifFollowing {
				err = trans.Insert(&ListenTo{ArtistId:artistId, UserId: userId, ListenCount: 1000, Timestamp:timestamp, IfFollowing: ifFollowing})
			} else {
				err = trans.Insert(&ListenTo{ArtistId:artistId, UserId: userId, ListenCount: int64(counter[artistId]), Timestamp:timestamp, IfFollowing: ifFollowing})
			}
		} else {
			_, err = dbmap.Exec("update ListenTo set listenCount = listenCount+?, timestamp = ? where artistId = ? and userId = ?", counter[artistId], timestamp, artistId, userId)
		}
		if err != nil {
			return err
		}
	}
	return trans.Commit()
}

// delete relation
func (d *MyDbMap) DeleteRelation(userId string, timestamp time.Time) error {
	_, err := d.DbMap.Exec("delete from ListenTo where userId = ? and timestamp < ? and ifFollowing = false", userId, timestamp)
	if err != nil {
		return err
	}
	return nil
}

// update the user's following artists
func (d *MyDbMap) UpdateFollowingRelation(artists []ArtistInfo, userId string, timestamp time.Time) error {
	fmt.Println(timestamp)
	trans, err := d.DbMap.Begin()
	if err != nil {
		return err
	}

	for _, artist := range artists {
		artistId := artist.ArtistId
		count, err := trans.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ? and ifFollowing = true", artistId, userId)
		if err != nil {
			return err
		}
		if count == 0 {
			if err := trans.Insert(&ListenTo{ArtistId:artistId, UserId: userId, ListenCount: 1000, Timestamp:timestamp, IfFollowing: true}); err != nil {
				return err
			}
		} else {
			// if _, err := trans.Update(&ListenTo{ArtistId:artistId, UserId: userId, Timestamp:timestamp, IfFollowing: true}); err != nil {
			// 	return err
			// }
			if _, err := trans.Exec("update ListenTo set timestamp = ? where artistId = ? and userId = ?", timestamp, artistId, userId); err != nil {
				return err
			}
		}
	}

	if err := trans.Commit(); err != nil {
		return err
	}

	return nil
}

// delete following relation if the user unfollowrd artists
func (d *MyDbMap) DeleteFollowingRelations(userId string, timestamp time.Time) error {
	fmt.Println(timestamp)
	if _, err := d.DbMap.Exec("delete from ListenTo where userId = ? and timestamp < ? and ifFollowing = true", userId, timestamp); err != nil {
		return err
	}

	return nil
}

func (d *MyDbMap) DeleteRelationFromRequest(userId string, artistIds []string) error {
	trans, err := d.DbMap.Begin()
	if err != nil {
		return err
	}

	for _, artistId := range artistIds {
		if _, err := d.DbMap.Exec("delete from ListenTo where userId = ? and artistId = ?", userId, artistId); err != nil {
			return err
		}
	}

	if err := trans.Commit(); err != nil {
		return err
	}

	return nil
}
