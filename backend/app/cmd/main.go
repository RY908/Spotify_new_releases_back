package main

import (
	"database/sql"
	"fmt"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/config"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/handlers"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao"
	_ "github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/dao/mysql"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/models/v2.0/schema"
	"github.com/RY908/Spotify_new_releases_back/backend/app/internal/usecase"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	config := config.LoadConfig()

	dbmap, err := newDB(config.DBConfig)
	if err != nil {
		panic(err)
	}
	dbManager, err := dao.NewDBManager(dbmap)
	if err != nil {
		panic(err)
	}

	createPlaylistUsecase := usecase.NewCreatePlaylistUsecase(dbManager, config.SpotifyConfig)
	getSettingUsecase := usecase.NewGetSettingUsecase(dbManager, config.SpotifyConfig)
	editSettingUsecase := usecase.NewEditSettingUsecase(dbManager, config.SpotifyConfig)
	deleteListeningHistoryUsecase := usecase.NewDeleteListeningHistoryUsecase(dbManager, config.SpotifyConfig)
	getArtistsByUserIDUsecase := usecase.NewGetArtistsByUserIDUsecase(dbManager, config.SpotifyConfig)

	s := newServer(createPlaylistUsecase,
		getSettingUsecase,
		editSettingUsecase,
		deleteListeningHistoryUsecase,
		getArtistsByUserIDUsecase)
	s.Start(":9990")
}

func newServer(
	createPlaylistUsecase *usecase.CreatePlaylistUsecase,
	getSettingUsecase *usecase.GetSettingUsecase,
	editSettingUsecase *usecase.EditSettingUsecase,
	deleteListeningHistoryUsecase *usecase.DeleteListeningHistoryUsecase,
	getArtistsByUserIDUsecase *usecase.GetArtistsByUserIDUsecase) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://labstack.com", "https://labstack.net"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	loginHandler := handlers.NewLoginHandler(createPlaylistUsecase)
	userHandler := handlers.NewUserHandler(deleteListeningHistoryUsecase, getArtistsByUserIDUsecase)
	settingHandler := handlers.NewSettinHandler(getSettingUsecase, editSettingUsecase)

	e.GET("/login", loginHandler.Login)
	e.GET("/callback", loginHandler.Callback)
	e.GET("/user", userHandler.GetArtists)
	e.POST("/delete", userHandler.DeleteArtists)
	e.GET("/setting", settingHandler.GetSettings)
	e.POST("/setting/save", settingHandler.EditSettings)

	return e
}

func newDB(config config.DBConfig) (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbmap.AddTableWithName(schema.Artist{}, "Artist").SetKeys(false, "ID")
	dbmap.AddTableWithName(schema.ListeningHistory{}, "ListenTo").SetKeys(true, "ID")
	dbmap.AddTableWithName(schema.User{}, "User").SetKeys(false, "ID")

	return dbmap, nil
}
