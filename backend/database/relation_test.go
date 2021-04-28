package database

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
	"testing"
	"time"
)

func TestInsertRelations(t *testing.T) {
	dbmap, err := DatabaseTestInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	// currentTime := time.Now().UTC()
	currentTime := time.Date(2021, time.March, 25, 0, 0, 0, 0, time.UTC)

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
			ArtistInfo{
				ArtistId: "existing_artistId4",
				Name:     "existing_artist4",
				Url:      "existing_url4",
				IconUrl:  "existing_iconUrl4",
			},
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
		map[string]int{"existing_artistId1": 1, "existing_artistId2": 2},
		"existing_user",
		currentTime,
		false,
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
			ArtistInfo{
				ArtistId: "existing_artistId4",
				Name:     "existing_artist4",
				Url:      "existing_url4",
				IconUrl:  "existing_iconUrl4",
			},
		},
		map[string]int{"existing_artistId3": 1, "existing_artistId4": 2},
		"existing_user",
		currentTime,
		true,
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []ListenTo
		wantErr error
	}{
		{
			name: "able to get relation",
			want: []ListenTo{
				ListenTo{
					ListenId:    1,
					ArtistId:    "existing_artistId1",
					UserId:      "existing_user",
					ListenCount: 1,
					Timestamp:   currentTime,
					IfFollowing: false,
				},
				ListenTo{
					ListenId:    2,
					ArtistId:    "existing_artistId2",
					UserId:      "existing_user",
					ListenCount: 2,
					Timestamp:   currentTime,
					IfFollowing: false,
				},
				ListenTo{
					ListenId:    3,
					ArtistId:    "existing_artistId3",
					UserId:      "existing_user",
					ListenCount: 1000,
					Timestamp:   currentTime,
					IfFollowing: true,
				},
				ListenTo{
					ListenId:    4,
					ArtistId:    "existing_artistId4",
					UserId:      "existing_user",
					ListenCount: 1000,
					Timestamp:   currentTime,
					IfFollowing: true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbmap.GetAllRelations()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ExistUser() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDeleteRelation(t *testing.T) {
	dbmap, err := DatabaseTestInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	// currentTime := time.Now().UTC()
	currentTime := time.Date(2021, time.March, 25, 0, 0, 0, 0, time.UTC)
	afterTime := currentTime.Add(1 * time.Hour)

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
			ArtistInfo{
				ArtistId: "existing_artistId4",
				Name:     "existing_artist4",
				Url:      "existing_url4",
				IconUrl:  "existing_iconUrl4",
			},
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
		},
		map[string]int{"existing_artistId1": 1},
		"existing_user",
		currentTime,
		false,
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.InsertRelations(
		[]ArtistInfo{
			ArtistInfo{
				ArtistId: "existing_artistId2",
				Name:     "existing_artist2",
				Url:      "existing_url2",
				IconUrl:  "existing_iconUrl2",
			},
		},
		map[string]int{"existing_artistId2": 2},
		"existing_user",
		afterTime,
		false,
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
		map[string]int{"existing_artistId3": 1},
		"existing_user",
		currentTime,
		true,
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.InsertRelations(
		[]ArtistInfo{
			ArtistInfo{
				ArtistId: "existing_artistId4",
				Name:     "existing_artist4",
				Url:      "existing_url4",
				IconUrl:  "existing_iconUrl4",
			},
		},
		map[string]int{"existing_artistId4": 2},
		"existing_user",
		afterTime,
		true,
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.DeleteRelation(
		"existing_user",
		afterTime,
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []ListenTo
		wantErr error
	}{
		{
			name: "able to get relation",
			want: []ListenTo{
				ListenTo{
					ListenId:    2,
					ArtistId:    "existing_artistId2",
					UserId:      "existing_user",
					ListenCount: 2,
					Timestamp:   afterTime,
					IfFollowing: false,
				},
				ListenTo{
					ListenId:    3,
					ArtistId:    "existing_artistId3",
					UserId:      "existing_user",
					ListenCount: 1000,
					Timestamp:   currentTime,
					IfFollowing: true,
				},
				ListenTo{
					ListenId:    4,
					ArtistId:    "existing_artistId4",
					UserId:      "existing_user",
					ListenCount: 1000,
					Timestamp:   afterTime,
					IfFollowing: true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbmap.GetAllRelations()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ExistUser() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDeleteRelationFromRequest(t *testing.T) {
	dbmap, err := DatabaseTestInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	// currentTime := time.Now().UTC()
	currentTime := time.Date(2021, time.March, 25, 0, 0, 0, 0, time.UTC)

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
			ArtistInfo{
				ArtistId: "existing_artistId4",
				Name:     "existing_artist4",
				Url:      "existing_url4",
				IconUrl:  "existing_iconUrl4",
			},
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
		map[string]int{"existing_artistId1": 1, "existing_artistId2": 2},
		"existing_user",
		currentTime,
		false,
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
			ArtistInfo{
				ArtistId: "existing_artistId4",
				Name:     "existing_artist4",
				Url:      "existing_url4",
				IconUrl:  "existing_iconUrl4",
			},
		},
		map[string]int{"existing_artistId3": 1, "existing_artistId4": 2},
		"existing_user",
		currentTime,
		true,
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.DeleteRelationFromRequest(
		"existing_user",
		[]string{
			"existing_artistId1",
			"existing_artistId3",
		},
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []ListenTo
		wantErr error
	}{
		{
			name: "able to get relation",
			want: []ListenTo{
				ListenTo{
					ListenId:    2,
					ArtistId:    "existing_artistId2",
					UserId:      "existing_user",
					ListenCount: 2,
					Timestamp:   currentTime,
					IfFollowing: false,
				},
				ListenTo{
					ListenId:    4,
					ArtistId:    "existing_artistId4",
					UserId:      "existing_user",
					ListenCount: 1000,
					Timestamp:   currentTime,
					IfFollowing: true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbmap.GetAllRelations()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ExistUser() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}

}
