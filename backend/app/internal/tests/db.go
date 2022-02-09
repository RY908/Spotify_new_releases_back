package tests

import (
	"database/sql"
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao/mysql"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

var (
	sqlPathTest = os.Getenv("SQL_PATH_TEST")
)

func NewTestDB() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", sqlPathTest)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbmap.AddTableWithName(schema.Artist{}, "Artist").SetKeys(false, "ID")
	dbmap.AddTableWithName(schema.ListeningHistory{}, "ListenTo").SetKeys(true, "ID")
	dbmap.AddTableWithName(schema.User{}, "User").SetKeys(false, "ID")
	//defer dbmap.Db.Close()
	//defer db.Close()
	//defer dbmap.Db.Close()

	return dbmap, nil
}

func TruncateTable(t *testing.T, d *gorp.DbMap) {
	t.Helper()

	if _, err := d.Exec("set foreign_key_checks = 0"); err != nil {
		t.Fatal(err)
	}
	if err := d.TruncateTables(); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Exec("set foreign_key_checks = 1"); err != nil {
		t.Fatal(err)
	}
}

func NewTestDBManager(db *gorp.DbMap) (dao.Factory, error) {
	return &mysql.DB{db}, nil
}
