package mysql_test

import (
	"errors"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/tests"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func Test_GetArtistsByUserID(t *testing.T) {
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
	userArtistsDao := dbManager.UserArtistsDAO()

	if err := insertArtists(artistDao, []*schema.Artist{testArtist1, testArtist2}); err != nil {
		t.Fatal(err)
	}
	if err := insertUsers(userDao, []*schema.User{testUser1, testUser2}); err != nil {
		t.Fatal(err)
	}
	if err := insertHistories(listeningHistoryDao, []*schema.ListeningHistory{testHisotry1, testHisotry2, testHisotry3, testHisotry4}); err != nil {
		t.Fatal(err)
	}

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*schema.UserArtist
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				userID: testUser1.ID,
			},
			want: []*schema.UserArtist{
				{
					ID:          testArtist1.ID,
					Name:        testArtist1.Name,
					Url:         testArtist1.Url,
					IconUrl:     testArtist1.IconUrl,
					IsFollowing: testHisotry1.IsFollowing,
				},
				{
					ID:          testArtist2.ID,
					Name:        testArtist2.Name,
					Url:         testArtist2.Url,
					IconUrl:     testArtist2.IconUrl,
					IsFollowing: testHisotry2.IsFollowing,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userArtistsDao.GetArtistsByUserID(tt.args.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetArtistsByUserID() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetArtistsByUserID() diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
