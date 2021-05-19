package database

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
)

func TestInsertArtist(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	if err := dbmap.InsertArtist(
		"existing_artist",
		"existing_artistId",
		"existing_url",
		"existing_iconUrl",
	); err != nil {
		t.Fatal(err)
	}
}

func TestGetArtistsFromUserId(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	currentTime := time.Now().UTC()

	if err := dbmap.InsertArtists(
		[]ArtistInfo{
			ArtistInfo{
				ArtistId: "existing_artistId1",
				Name:     "existing_artist1",
				Url:      "existing_url1",
				IconUrl:  "existing_iconUrl1",
			},
			ArtistInfo{
				ArtistId: "existing_artistId2",
				Name:     "existing_artist2",
				Url:      "existing_url2",
				IconUrl:  "existing_iconUrl2",
			},
			ArtistInfo{
				ArtistId: "existing_artistId3",
				Name:     "existing_artist3",
				Url:      "existing_url3",
				IconUrl:  "existing_iconUrl3",
			},
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.InsertUser(
		"existing_user",
		"existing_playlist",
		&oauth2.Token{
			AccessToken:  "existing_accessToken",
			TokenType:    "existing_tokenType",
			RefreshToken: "existing_refreshToken",
			Expiry:       currentTime,
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.InsertRelations(
		[]ArtistInfo{
			ArtistInfo{
				ArtistId: "existing_artistId1",
				Name:     "existing_artist1",
				Url:      "existing_url1",
				IconUrl:  "existing_iconUrl1",
			},
			ArtistInfo{
				ArtistId: "existing_artistId2",
				Name:     "existing_artist2",
				Url:      "existing_url2",
				IconUrl:  "existing_iconUrl2",
			},
		},
		map[string]int{"existing_artistId1": 1, "existing_artistId": 2},
		"existing_user",
		currentTime,
		true,
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.InsertRelations(
		[]ArtistInfo{
			ArtistInfo{
				ArtistId: "existing_artistId3",
				Name:     "existing_artist3",
				Url:      "existing_url3",
				IconUrl:  "existing_iconUrl3",
			},
		},
		map[string]int{"existing_artistId3": 3},
		"existing_user",
		currentTime,
		false,
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		userId  string
		want    []ArtistRes
		wantErr error
	}{
		{
			name:   "able to get artists",
			userId: "existing_user",
			want: []ArtistRes{
				ArtistRes{
					ArtistId:    "existing_artistId1",
					Name:        "existing_artist1",
					Url:         "existing_url1",
					IconUrl:     "existing_iconUrl1",
					IfFollowing: true,
				},
				ArtistRes{
					ArtistId:    "existing_artistId2",
					Name:        "existing_artist2",
					Url:         "existing_url2",
					IconUrl:     "existing_iconUrl2",
					IfFollowing: true,
				},
				ArtistRes{
					ArtistId:    "existing_artistId3",
					Name:        "existing_artist3",
					Url:         "existing_url3",
					IconUrl:     "existing_iconUrl3",
					IfFollowing: false,
				},
			},
			wantErr: nil,
		},
		{
			name:    "able to return error for not existing user",
			userId:  "nonexisting_user",
			want:    []ArtistRes{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbmap.GetArtistsFromUserId(tt.userId)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetArtistsFromUserId() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetArtistsFromUserId() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
