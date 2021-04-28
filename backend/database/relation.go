package database

import (
	"fmt"
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	"time"
)

type ListenTo struct {
	ListenId    int64     `db:listenId`
	ArtistId    string    `db:artistId`
	UserId      string    `db:userId`
	ListenCount int64     `db:listenCount`
	Timestamp   time.Time `db:"timestamp"`
	IfFollowing bool      `db:"ifFollowing"`
}

func (d *MyDbMap) GetAllRelations() ([]ListenTo, error) {
	var listenTo []ListenTo
	if _, err := d.DbMap.Select(&listenTo, "select * from ListenTo"); err != nil {
		return nil, fmt.Errorf("get all relations: %w", err)
	}

	return listenTo, nil
}

// insert relation
func (d *MyDbMap) InsertRelation(artistId, userId string, timestamp time.Time, ifFollowing bool) error {
	count, err := d.DbMap.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ?", artistId, userId)
	if err != nil {
		return fmt.Errorf("insert relation select count: %w", err)
	}
	if count == 0 {
		err = d.DbMap.Insert(&ListenTo{ArtistId: artistId, UserId: userId, Timestamp: timestamp, IfFollowing: ifFollowing})
	}
	if err != nil {
		return fmt.Errorf("insert relation: %w", err)
	}
	return nil
}

// insert multiple relations
func (d *MyDbMap) InsertRelations(artists []ArtistInfo, counter map[string]int, userId string, timestamp time.Time, ifFollowing bool) error {
	trans, err := d.DbMap.Begin()
	if err != nil {
		return fmt.Errorf("insert relations transaction begin: %w", err)
	}

	for _, artist := range artists {
		artistId := artist.ArtistId
		count, err := trans.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ?", artistId, userId)
		if err != nil {
			return fmt.Errorf("insert relation select count: %w", err)
		}
		if count == 0 {
			if ifFollowing {
				err = trans.Insert(&ListenTo{ArtistId: artistId, UserId: userId, ListenCount: 1000, Timestamp: timestamp, IfFollowing: ifFollowing})
			} else {
				err = trans.Insert(&ListenTo{ArtistId: artistId, UserId: userId, ListenCount: int64(counter[artistId]), Timestamp: timestamp, IfFollowing: ifFollowing})
			}
		} else {
			_, err = d.DbMap.Exec("update ListenTo set listenCount = listenCount+?, timestamp = ? where artistId = ? and userId = ?", counter[artistId], timestamp, artistId, userId)
		}
		if err != nil {
			return fmt.Errorf("insert relations: %w", err)
		}
	}

	if err := trans.Commit(); err != nil {
		return fmt.Errorf("insert relation transaction commit: %w", err)
	}

	return nil
}

// delete relation
func (d *MyDbMap) DeleteRelation(userId string, timestamp time.Time) error {
	_, err := d.DbMap.Exec("delete from ListenTo where userId = ? and timestamp < ? and ifFollowing = false", userId, timestamp)
	if err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}
	return nil
}

// update the user's following artists
func (d *MyDbMap) UpdateFollowingRelation(artists []ArtistInfo, userId string, timestamp time.Time) error {
	trans, err := d.DbMap.Begin()
	if err != nil {
		return fmt.Errorf("update following relation transaction begin: %w", err)
	}

	for _, artist := range artists {
		artistId := artist.ArtistId
		count, err := trans.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ? and ifFollowing = true", artistId, userId)
		if err != nil {
			return fmt.Errorf("update following relation select count: %w", err)
		}
		if count == 0 {
			if err := trans.Insert(&ListenTo{ArtistId: artistId, UserId: userId, ListenCount: 1000, Timestamp: timestamp, IfFollowing: true}); err != nil {
				return fmt.Errorf("update following relation insert: %w", err)
			}
		} else {
			// if _, err := trans.Update(&ListenTo{ArtistId:artistId, UserId: userId, Timestamp:timestamp, IfFollowing: true}); err != nil {
			// 	return err
			// }
			if _, err := trans.Exec("update ListenTo set timestamp = ? where artistId = ? and userId = ?", timestamp, artistId, userId); err != nil {
				return fmt.Errorf("update following relation update: %w", err)
			}
		}
	}

	if err := trans.Commit(); err != nil {
		return fmt.Errorf("update following relation transaction commit: %w", err)
	}

	return nil
}

// delete following relation if the user unfollowrd artists
func (d *MyDbMap) DeleteFollowingRelations(userId string, timestamp time.Time) error {
	if _, err := d.DbMap.Exec("delete from ListenTo where userId = ? and timestamp < ? and ifFollowing = true", userId, timestamp); err != nil {
		return fmt.Errorf("delete following relation: %w", err)
	}

	return nil
}

func (d *MyDbMap) DeleteRelationFromRequest(userId string, artistIds []string) error {
	trans, err := d.DbMap.Begin()
	if err != nil {
		return fmt.Errorf("delete following relation from request transaction begin: %w", err)
	}

	for _, artistId := range artistIds {
		if _, err := d.DbMap.Exec("delete from ListenTo where userId = ? and artistId = ?", userId, artistId); err != nil {
			return fmt.Errorf("delete following relation from request delete: %w", err)
		}
	}

	if err := trans.Commit(); err != nil {
		return fmt.Errorf("delete following relation from request transaction commit: %w", err)
	}

	return nil
}
