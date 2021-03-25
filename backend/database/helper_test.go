package database

import (
	"testing"

	// "github.com/go-gorp/gorp/v3"
)

func truncateTable(t *testing.T, d *MyDbMap) {
	t.Helper()

	if _, err := d.DbMap.Exec("set foreign_key_checks = 0"); err != nil {
		t.Fatal(err)
	}
	if err := d.DbMap.TruncateTables(); err != nil {
		t.Fatal(err)
	}
	if _, err := d.DbMap.Exec("set foreign_key_checks = 1"); err != nil {
		t.Fatal(err)
	}
}