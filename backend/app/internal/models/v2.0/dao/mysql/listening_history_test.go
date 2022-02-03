package mysql_test

import (
	"database/sql"
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/tests"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

var (
	currentTime  = time.Now()
	testHisotry1 = &schema.ListeningHistory{
		ID:          1,
		UserID:      "test_userID1",
		ArtistID:    "test_artistID1",
		ListenCount: 1000,
		Timestamp:   currentTime,
		IsFollowing: true,
	}
	testHisotry2 = &schema.ListeningHistory{
		ID:          2,
		UserID:      "test_userID1",
		ArtistID:    "test_artistID2",
		ListenCount: 1000,
		Timestamp:   currentTime,
		IsFollowing: true,
	}
	testHisotry3 = &schema.ListeningHistory{
		ID:          3,
		UserID:      "test_userID2",
		ArtistID:    "test_artistID2",
		ListenCount: 1000,
		Timestamp:   time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC),
		IsFollowing: false,
	}
	testHisotry4 = &schema.ListeningHistory{
		ID:          4,
		UserID:      "test_userID2",
		ArtistID:    "test_artistID1",
		ListenCount: 1000,
		Timestamp:   currentTime,
		IsFollowing: false,
	}
)

func Test_GetAllListeningHistory(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}
	artistDao := dbManager.ArtistDAO()
	listeningHistoryDao := dbManager.ListeningHistoryDAO()
	userDao := dbManager.UserDAO()

	if err := insertArtists(artistDao, []*schema.Artist{testArtist1}); err != nil {
		t.Fatal(err)
	}
	if err := insertUsers(userDao, []*schema.User{testUser1}); err != nil {
		t.Fatal(err)
	}
	if err := insertHistories(listeningHistoryDao, []*schema.ListeningHistory{testHisotry1}); err != nil {
		t.Fatal(err)
	}

	type args struct {
		artistID string
		userID   string
	}

	tests := []struct {
		name    string
		args    args
		want    *schema.ListeningHistory
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				artistID: "test_artistID1",
				userID:   "test_userID1",
			},
			want:    testHisotry1,
			wantErr: nil,
		},
		{
			name: "not ok: listening history does not exist",
			args: args{
				artistID: "test_artistID2",
				userID:   "test_userID1",
			},
			want:    nil,
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listeningHistoryDao.GetListeningHistory(tt.args.artistID, tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAllListeningHistory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.ListeningHistory{}, "Timestamp"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("GetAllListeningHistory() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_GetAllListeningHistories(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}
	artistDao := dbManager.ArtistDAO()
	listeningHistoryDao := dbManager.ListeningHistoryDAO()
	userDao := dbManager.UserDAO()

	if err := insertArtists(artistDao, []*schema.Artist{testArtist1, testArtist2}); err != nil {
		t.Fatal(err)
	}
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}
	if err := insertHistories(listeningHistoryDao, []*schema.ListeningHistory{testHisotry1, testHisotry2, testHisotry3}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []*schema.ListeningHistory
		wantErr error
	}{
		{
			name:    "ok",
			want:    []*schema.ListeningHistory{testHisotry1, testHisotry2, testHisotry3},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listeningHistoryDao.GetAllListeningHistory()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAllListeningHistories() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.ListeningHistory{}, "Timestamp"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("GetAllListeningHistories() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_UpdateHistory(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}
	artistDao := dbManager.ArtistDAO()
	listeningHistoryDao := dbManager.ListeningHistoryDAO()
	userDao := dbManager.UserDAO()

	if err := insertArtists(artistDao, []*schema.Artist{testArtist1}); err != nil {
		t.Fatal(err)
	}
	if err := insertUsers(userDao, []*schema.User{testUser1}); err != nil {
		t.Fatal(err)
	}
	if err := insertHistories(listeningHistoryDao, []*schema.ListeningHistory{testHisotry1}); err != nil {
		t.Fatal(err)
	}

	if err := listeningHistoryDao.UpdateHistory(testArtist1.ID, testUser1.ID, 10, true, time.Now()); err != nil {
		t.Fatal(err)
	}

	type args struct {
		artistID string
		userID   string
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.ListeningHistory
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				artistID: "test_artistID1",
				userID:   "test_userID1",
			},
			want: &schema.ListeningHistory{
				ID:          1,
				UserID:      "test_userID1",
				ArtistID:    "test_artistID1",
				ListenCount: 1010,
				Timestamp:   time.Now(),
				IsFollowing: true,
			},
			wantErr: nil,
		},
		{
			name: "not ok: listening history does not exist",
			args: args{
				artistID: "test_artistID2",
				userID:   "test_userID1",
			},
			want:    nil,
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listeningHistoryDao.GetListeningHistory(tt.args.artistID, tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateHistory() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.ListeningHistory{}, "Timestamp"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("UpdateHistory() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_DeleteHistoryByTimestamp(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}
	artistDao := dbManager.ArtistDAO()
	listeningHistoryDao := dbManager.ListeningHistoryDAO()
	userDao := dbManager.UserDAO()

	if err := insertArtists(artistDao, []*schema.Artist{testArtist1, testArtist2}); err != nil {
		t.Fatal(err)
	}
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}
	if err := insertHistories(listeningHistoryDao, []*schema.ListeningHistory{testHisotry1, testHisotry2, testHisotry3, testHisotry4}); err != nil {
		t.Fatal(err)
	}

	if err := listeningHistoryDao.DeleteHistoryByTimestamp(testUser1.ID, currentTime.Add(-1*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := listeningHistoryDao.DeleteHistoryByTimestamp(testUser2.ID, currentTime.Add(-1*time.Hour)); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []*schema.ListeningHistory
		wantErr error
	}{
		{
			name:    "ok",
			want:    []*schema.ListeningHistory{testHisotry1, testHisotry2, testHisotry4},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listeningHistoryDao.GetAllListeningHistory()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteHistoryByTimestamp() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.ListeningHistory{}, "Timestamp", "ID"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("DeleteHistoryByTimestamp() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}

}

func Test_DeleteHistoryByArtistId(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}
	artistDao := dbManager.ArtistDAO()
	listeningHistoryDao := dbManager.ListeningHistoryDAO()
	userDao := dbManager.UserDAO()

	if err := insertArtists(artistDao, []*schema.Artist{testArtist1, testArtist2}); err != nil {
		t.Fatal(err)
	}
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}
	if err := insertHistories(listeningHistoryDao, []*schema.ListeningHistory{testHisotry1, testHisotry2, testHisotry3, testHisotry4}); err != nil {
		t.Fatal(err)
	}

	if err := listeningHistoryDao.DeleteHistoryByArtistId(testUser1.ID, testArtist1.ID); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []*schema.ListeningHistory
		wantErr error
	}{
		{
			name:    "ok",
			want:    []*schema.ListeningHistory{testHisotry2, testHisotry3, testHisotry4},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listeningHistoryDao.GetAllListeningHistory()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteHistoryByTimestamp() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.ListeningHistory{}, "Timestamp", "ID"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("DeleteHistoryByTimestamp() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func insertHistories(dao dao.ListeningHistory, listeningHistories []*schema.ListeningHistory) error {
	for _, history := range listeningHistories {
		if err := dao.InsertHistory(history); err != nil {
			return err
		}
	}
	return nil
}
