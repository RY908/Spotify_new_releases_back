package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-gorp/gorp"
)

type MyDbMap struct {
	DbMap *gorp.DbMap
}

func DatabaseInit(dbmap *gorp.DbMap) *MyDbMap {
	dbmap.AddTableWithName(ArtistInfo{}, "Artist").SetKeys(false, "ArtistId")
	dbmap.AddTableWithName(ListenTo{}, "ListenTo").SetKeys(true, "ListenId")
	dbmap.AddTableWithName(UserInfo{}, "User").SetKeys(false, "UserId")
	//defer dbmap.Db.Close()

	return &MyDbMap{DbMap: dbmap}
}