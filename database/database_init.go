package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-gorp/gorp"
	"os"
	"database/sql"
)

var (
	sqlPath = os.Getenv("SQL_PATH")
	db, _ = sql.Open("mysql", sqlPath)
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
)

type MyDbMap struct {
	DbMap *gorp.DbMap
}

func DatabaseInit() *MyDbMap {
	dbmap.AddTableWithName(ArtistInfo{}, "Artist").SetKeys(false, "ArtistId")
	dbmap.AddTableWithName(ListenTo{}, "ListenTo").SetKeys(true, "ListenId")
	dbmap.AddTableWithName(UserInfo{}, "User").SetKeys(false, "UserId")
	//defer dbmap.Db.Close()
	//defer db.Close()
	//defer dbmap.Db.Close()

	return &MyDbMap{DbMap: dbmap}
}