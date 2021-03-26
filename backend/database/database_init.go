package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-gorp/gorp"
	"os"
	"database/sql"
	"fmt"
)

var (
	sqlPath = os.Getenv("SQL_PATH")
	sqlPathTest = os.Getenv("SQL_PATH_TEST")
)

type MyDbMap struct {
	DbMap *gorp.DbMap
}

func DatabaseInit() (*MyDbMap, error) {
	db, err := sql.Open("mysql", sqlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbmap.AddTableWithName(ArtistInfo{}, "Artist").SetKeys(false, "ArtistId")
	dbmap.AddTableWithName(ListenTo{}, "ListenTo").SetKeys(true, "ListenId")
	dbmap.AddTableWithName(UserInfo{}, "User").SetKeys(false, "UserId")
	//defer dbmap.Db.Close()
	//defer db.Close()
	//defer dbmap.Db.Close()

	return &MyDbMap{DbMap: dbmap}, nil
}

func DatabaseTestInit() (*MyDbMap, error) {
	fmt.Println(sqlPathTest)
	db, err := sql.Open("mysql", sqlPathTest)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbmap.AddTableWithName(ArtistInfo{}, "Artist").SetKeys(false, "ArtistId")
	dbmap.AddTableWithName(ListenTo{}, "ListenTo").SetKeys(true, "ListenId")
	dbmap.AddTableWithName(UserInfo{}, "User").SetKeys(false, "UserId")
	//defer dbmap.Db.Close()
	//defer db.Close()
	//defer dbmap.Db.Close()

	return &MyDbMap{DbMap: dbmap}, nil
}