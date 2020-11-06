package database

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/go-gorp/gorp"
	"fmt"
	"time"
)

type ListenTo struct {
	ListenId	int64 	`db:listenId`
	ArtistId	string 	`db:artistId`
	UserId		string 	`db:userId`
	Timestamp	time.Time	`db:"timestamp"`
}

func (d *MyDbMap) InsertRelation(artistId, userId string, timestamp time.Time) error {
	count, err := d.DbMap.SelectInt("select count(*) from ListenTo where artistId = ? and userId = ?", artistId, userId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if count == 0 {
		err = d.DbMap.Insert(&ListenTo{ArtistId:artistId, UserId: userId, Timestamp:timestamp})
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (d *MyDbMap) DeleteRelation(userId string, timestamp time.Time) error {
	res, err := d.DbMap.Exec("delete from ListenTo where userId = ? and timestamp <> ?", userId, timestamp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}