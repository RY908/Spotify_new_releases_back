package tests

import (
	"database/sql"
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
	"os"
)

var (
	sqlPathTest = os.Getenv("SQL_PATH_TEST")
)

func DatabaseTestInit() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", sqlPathTest)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbmap.AddTableWithName(schema.Artist{}, "Artist").SetKeys(false, "ArtistId")
	dbmap.AddTableWithName(schema.ListeningHistory{}, "ListenTo").SetKeys(true, "ListenId")
	dbmap.AddTableWithName(schema.User{}, "User").SetKeys(false, "UserId")
	//defer dbmap.Db.Close()
	//defer db.Close()
	//defer dbmap.Db.Close()

	return dbmap, nil
}
