package mysql_test

import (
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
	testUser1 = &schema.User{
		ID:            "test_userID1",
		AccessToken:   "test_user_accessToken1",
		TokenType:     "test_user_tokenType1",
		RefreshToken:  "test_user_refreshToken1",
		Expiry:        time.Now().Add(24 * time.Hour),
		PlaylistId:    "test_user_playlistID1",
		IfRemixAdd:    true,
		IfAcousticAdd: true,
	}
	testUser2 = &schema.User{
		ID:            "test_userID2",
		AccessToken:   "test_user_accessToken2",
		TokenType:     "test_user_tokenType2",
		RefreshToken:  "test_user_refreshToken2",
		Expiry:        time.Now().Add(24 * time.Hour),
		PlaylistId:    "test_user_playlistID2",
		IfRemixAdd:    false,
		IfAcousticAdd: false,
	}
)

func Test_GetUser(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}

	userDao := dbManager.UserDAO()
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.User
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				userID: testUser1.ID,
			},
			want:    testUser1,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDao.GetUser(tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.User{}, "Expiry", "ID"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("GetUser() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_GetAllUsers(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}

	userDao := dbManager.UserDAO()
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []*schema.User
		wantErr error
	}{
		{
			name:    "ok",
			want:    []*schema.User{testUser1, testUser2},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDao.GetAllUsers()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAllUsers() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.User{}, "Expiry", "ID"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("GetAllUsers() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_UpdateUser(t *testing.T) {
	db, err := tests.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}
	tests.TruncateTable(t, db)
	dbManager, err := dao.NewDBManager(db)
	if err != nil {
		t.Fatal(err)
	}

	userDao := dbManager.UserDAO()
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}

	testUser1.AccessToken = "updated_test_accessToken1"
	if err := userDao.UpdateUser(testUser1); err != nil {
		t.Fatal(err)
	}

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.User
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				userID: testUser1.ID,
			},
			want:    testUser1,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDao.GetUser(tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdateUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(schema.User{}, "Expiry", "ID"),
			}
			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("UpdateUser() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func insertUsers(dao dao.User, users []*schema.User) error {
	for _, user := range users {
		if err := dao.InsertUser(user); err != nil {
			return err
		}
	}
	return nil
}
