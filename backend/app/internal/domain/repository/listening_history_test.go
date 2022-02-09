package repository_test

import (
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/entity"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/domain/repository"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/tests"
	test_models "github.com/RY908/Spotify_new_releases_back/backend/app/internal/tests/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

func Test_InsertOrUpdateListeningHistory(t *testing.T) {
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
	listeningHistoryRepo := repository.NewListeningHistoryRepository()

	test_models.InsertUsers(userDao)
	test_models.InsertArtists(artistDao)
	test_models.InsertHistories(listeningHistoryDao)

	currentTime := time.Now()

	type args struct {
		listeningHistory entity.ListeningHistory
		count            int
		artistID         string
		userID           string
	}

	tests := []struct {
		name    string
		args    args
		want    *entity.ListeningHistory
		wantErr error
	}{
		{
			name: "ok: insert new listening history with no following",
			args: args{
				listeningHistory: entity.ListeningHistory{
					UserID:      "test_userID3",
					ArtistID:    "test_artistID3",
					ListenCount: 1,
					Timestamp:   currentTime,
					IsFollowing: false,
				},
				count:    1,
				artistID: "test_artistID3",
				userID:   "test_userID3",
			},
			want: &entity.ListeningHistory{
				UserID:      "test_userID3",
				ArtistID:    "test_artistID3",
				ListenCount: 1,
				Timestamp:   currentTime,
				IsFollowing: false,
			},
			wantErr: nil,
		},
		{
			name: "ok: insert new listening history with following",
			args: args{
				listeningHistory: entity.ListeningHistory{
					UserID:      "test_userID2",
					ArtistID:    "test_artistID3",
					ListenCount: 1,
					Timestamp:   currentTime,
					IsFollowing: true,
				},
				count:    1,
				artistID: "test_artistID3",
				userID:   "test_userID2",
			},
			want: &entity.ListeningHistory{
				UserID:      "test_userID2",
				ArtistID:    "test_artistID3",
				ListenCount: 1000,
				Timestamp:   currentTime,
				IsFollowing: true,
			},
			wantErr: nil,
		},
		{
			name: "ok: update listening history",
			args: args{
				listeningHistory: entity.ListeningHistory{
					UserID:      "test_userID1",
					ArtistID:    "test_artistID1",
					ListenCount: 100,
					Timestamp:   currentTime,
					IsFollowing: false,
				},
				count:    100,
				artistID: "test_artistID1",
				userID:   "test_userID1",
			},
			want: &entity.ListeningHistory{
				UserID:      "test_userID1",
				ArtistID:    "test_artistID1",
				ListenCount: 1100,
				Timestamp:   currentTime,
				IsFollowing: true,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := listeningHistoryRepo.InsertOrUpdateListeningHistory(dbManager, tt.args.listeningHistory, tt.args.count)
			if err != nil {
				t.Fatal(err)
			}
			got, err := listeningHistoryDao.GetListeningHistory(tt.args.artistID, tt.args.userID)
			listeningHistory := entity.NewListeningHistory(got)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAllListeningHistory() error = %v, wantErr = %v", err, tt.wantErr)
			}

			opts := []cmp.Option{
				cmpopts.IgnoreFields(entity.ListeningHistory{}, "Timestamp", "ID"),
			}
			if !cmp.Equal(listeningHistory, tt.want, opts...) {
				t.Errorf("GetAllListeningHistory() diff = %v", cmp.Diff(listeningHistory, tt.want))
			}
		})
	}
}

func Test_UpdateIsFollowing(t *testing.T) {
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
	listeningHistoryRepo := repository.NewListeningHistoryRepository()

	test_models.InsertUsers(userDao)
	test_models.InsertArtists(artistDao)
	test_models.InsertHistories(listeningHistoryDao)

	currentTime := time.Now()

	type args struct {
		artistID string
		userID   string
	}

	tests := []struct {
		name    string
		args    args
		want    *entity.ListeningHistory
		wantErr error
	}{
		{
			name: "ok: update listening history which does not have following relation",
			args: args{
				artistID: "test_artistID2",
				userID:   "test_userID2",
			},
			want: &entity.ListeningHistory{
				UserID:      "test_userID2",
				ArtistID:    "test_artistID2",
				ListenCount: 1001,
				Timestamp:   currentTime,
				IsFollowing: true,
			},
			wantErr: nil,
		},
		{
			name: "ok: update listening history which alreadly has following relation",
			args: args{
				artistID: "test_artistID1",
				userID:   "test_userID1",
			},
			want: &entity.ListeningHistory{
				UserID:      "test_userID1",
				ArtistID:    "test_artistID1",
				ListenCount: 1000,
				Timestamp:   currentTime,
				IsFollowing: true,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := listeningHistoryRepo.UpdateIsFollowing(dbManager, tt.args.artistID, tt.args.userID, currentTime)
			if err != nil {
				t.Fatal(err)
			}
			got, err := listeningHistoryDao.GetListeningHistory(tt.args.artistID, tt.args.userID)
			listeningHistory := entity.NewListeningHistory(got)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAllListeningHistory() error = %v, wantErr = %v", err, tt.wantErr)
			}

			opts := []cmp.Option{
				cmpopts.IgnoreFields(entity.ListeningHistory{}, "Timestamp", "ID"),
			}
			if !cmp.Equal(listeningHistory, tt.want, opts...) {
				t.Errorf("GetAllListeningHistory() diff = %v", cmp.Diff(listeningHistory, tt.want))
			}
		})
	}
}
