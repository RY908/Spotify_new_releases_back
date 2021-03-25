package database

import (
	"errors"
	"testing"
	"time"
	"golang.org/x/oauth2"
	"github.com/google/go-cmp/cmp"
)

func TestInsertUser(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	currentTime := time.Now().UTC()
	
	if err := dbmap.InsertUser(
		"existing_user1", 
		"existing_playlist1",
		&oauth2.Token{
			AccessToken:	"existing_accessToken1",
			TokenType:		"existing_tokenType1",
			RefreshToken:	"existing_refreshToken1",
			Expiry:			currentTime, 
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestExistUser(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	currentTime := time.Now().UTC()
	
	if err := dbmap.InsertUser(
		"existing_user", 
		"existing_playlist",
		&oauth2.Token{
			AccessToken:	"existing_accessToken",
			TokenType:		"existing_tokenType",
			RefreshToken:	"existing_refreshToken",
			Expiry:			currentTime, 
		},
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name 		string
		userId 		string
		wantExists	bool
		want		UserInfo
		wantErr 	error
	} {
		{
			name: "able to get existing user", 
			userId: "existing_user",
			wantExists: true,
			want: UserInfo{
				UserId:			"existing_user",
				AccessToken: 	"existing_accessToken",
				TokenType: 		"existing_tokenType",
				RefreshToken: 	"existing_refreshToken",
				Expiry: 		currentTime, 
				PlaylistId: 	"existing_playlistId",
				IfRemixAdd: 	true,
				IfAcousticAdd: 	true,
			},
			wantErr: nil,
		},
		// {
		// 	name: "able to return false to nonexisting user", 
		// 	userId: "nonexisting_user1",
		// 	wantExists: false,
		// 	want: UserInfo{},
		// 	wantErr: nil, 
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, got, err := dbmap.ExistUser(tt.userId)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ExistUser() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotExists, tt.wantExists) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(gotExists, tt.wantExists))
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}

}

func TestGetAllUsers(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	currentTime := time.Now().UTC()
	
	if err := dbmap.InsertUser(
		"existing_user1", 
		"existing_playlist1",
		&oauth2.Token{
			AccessToken:	"existing_accessToken1",
			TokenType:		"existing_tokenType1",
			RefreshToken:	"existing_refreshToken1",
			Expiry:			currentTime, 
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.InsertUser(
		"existing_user2", 
		"existing_playlist2",
		&oauth2.Token{
			AccessToken:	"existing_accessToken2",
			TokenType:		"existing_tokenType2",
			RefreshToken:	"existing_refreshToken2",
			Expiry:			currentTime, 
		},
	); err != nil {
		t.Fatal(err)
	}


	tests := []struct {
		name 		string
		userId 		string
		want		[]UserInfo
		wantErr 	error
	} {
		{
			name: "able to get all existing user", 
			userId: "existing_user",
			want: []UserInfo{
				UserInfo{	
					UserId:			"existing_user1",
					AccessToken: 	"existing_accessToken1",
					TokenType: 		"existing_tokenType1",
					RefreshToken: 	"existing_refreshToken1",
					Expiry: 		currentTime, 
					PlaylistId: 	"existing_playlistId1",
					IfRemixAdd: 	true,
					IfAcousticAdd: 	true,
				},
				UserInfo{	
					UserId:			"existing_user2",
					AccessToken: 	"existing_accessToken2",
					TokenType: 		"existing_tokenType2",
					RefreshToken: 	"existing_refreshToken2",
					Expiry: 		currentTime, 
					PlaylistId: 	"existing_playlistId2",
					IfRemixAdd: 	true,
					IfAcousticAdd: 	true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbmap.GetAllUsers()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAllUsers() error: %v, wantErr: %v", err, tt.wantErr)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetAllUsers() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}

}

func TestUpdateUser(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	currentTime := time.Now().UTC()
	
	if err := dbmap.InsertUser(
		"existing_user", 
		"existing_playlist",
		&oauth2.Token{
			AccessToken:	"existing_accessToken",
			TokenType:		"existing_tokenType",
			RefreshToken:	"existing_refreshToken",
			Expiry:			currentTime, 
		},
	); err != nil {
		t.Fatal(err)
	}

	updatedTime := currentTime.Add(1 * time.Hour)

	if err := dbmap.UpdateUser(
		"existing_user",
		"existing_playlist",
		&oauth2.Token{
			AccessToken:	"updaeted_accessToken",
			TokenType:		"updated_tokenType",
			RefreshToken:	"updated_refreshToken",
			Expiry:			updatedTime, 
		},
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name 		string
		userId 		string
		wantExists	bool
		want		UserInfo
		wantErr 	error
	} {
		{
			name: "able to get updated user", 
			userId: "existing_user",
			wantExists: true,
			want: UserInfo{
				UserId:			"existing_user",
				AccessToken: 	"updated_accessToken",
				TokenType: 		"updated_tokenType",
				RefreshToken: 	"updated_refreshToken",
				Expiry: 		updatedTime, 
				PlaylistId: 	"updated_playlistId",
				IfRemixAdd: 	true,
				IfAcousticAdd: 	true,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, got, err := dbmap.ExistUser(tt.userId)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ExistUser() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotExists, tt.wantExists) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(gotExists, tt.wantExists))
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}

}

func TestUpdateIfAdd(t *testing.T) {
	dbmap, err := DatabaseInit()
	if err != nil {
		t.Fatal(err)
	}
	truncateTable(t, dbmap)

	currentTime := time.Now().UTC()
	if err := dbmap.InsertUser(
		"existing_user", 
		"existing_playlist",
		&oauth2.Token{
			AccessToken:	"existing_accessToken",
			TokenType:		"existing_tokenType",
			RefreshToken:	"existing_refreshToken",
			Expiry:			currentTime, 
		},
	); err != nil {
		t.Fatal(err)
	}

	if err := dbmap.UpdateIfAdd(
		"existing_user", 
		false,
		false,
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name 		string
		userId 		string
		wantExists	bool
		want		UserInfo
		wantErr 	error
	} {
		{
			name: "able to get updated user", 
			userId: "existing_user",
			wantExists: true,
			want: UserInfo{
				UserId:			"existing_user",
				AccessToken: 	"existing_accessToken",
				TokenType: 		"existing_tokenType",
				RefreshToken: 	"existing_refreshToken",
				Expiry: 		currentTime, 
				PlaylistId: 	"existing_playlistId",
				IfRemixAdd: 	false,
				IfAcousticAdd: 	false,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, got, err := dbmap.ExistUser(tt.userId)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ExistUser() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotExists, tt.wantExists) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(gotExists, tt.wantExists))
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("ExistUser() diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}