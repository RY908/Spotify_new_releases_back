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
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
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

	appLogger, dbLogger, cronLogger := newLoggers()

	// usecases for handlers
	createPlaylistUsecase := usecase.NewCreatePlaylistUsecase(dbManager, config.SpotifyConfig)
	getSettingUsecase := usecase.NewGetSettingUsecase(dbManager, config.SpotifyConfig)
	editSettingUsecase := usecase.NewEditSettingUsecase(dbManager, config.SpotifyConfig)
	deleteListeningHistoryUsecase := usecase.NewDeleteListeningHistoryUsecase(dbManager, config.SpotifyConfig)
	getArtistsByUserIDUsecase := usecase.NewGetArtistsByUserIDUsecase(dbManager, config.SpotifyConfig)

	// usecases for cron
	updateListeningHistoryUsecase := usecase.NewUpdateListeningHistoryUsecase(dbManager, config.SpotifyConfig)
	updatePlaylistUsecase := usecase.NewUpdatePlaylistUsecase(dbManager, config.SpotifyConfig)
	updateFollowingArtistsUsecase := usecase.NewUpdateFollowingArtistsUsecase(dbManager, config.SpotifyConfig)

	s := newServer(
		appLogger,
		config.CallbackConfig,
		createPlaylistUsecase,
		getSettingUsecase,
		editSettingUsecase,
		deleteListeningHistoryUsecase,
		getArtistsByUserIDUsecase)
	s.Start(":9990")

	c := newCron(
		cronLogger,
		updateListeningHistoryUsecase,
		updatePlaylistUsecase,
		updateFollowingArtistsUsecase)
	c.Start()
}

func newLoggers() (*log.Logger, *log.Logger, *log.Logger) {
	appLogger := log.New(os.Stdout, "[APP]", log.LstdFlags|log.LUTC)
	dbLogger := log.New(os.Stdout, "[DB]", log.LstdFlags|log.LUTC)
	cronLogger := log.New(os.Stdout, "[CRON]", log.LstdFlags|log.LUTC)
	return appLogger, dbLogger, cronLogger
}

func newServer(
	appLogger *log.Logger,
	config *config.CallbackConfig,
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

	loginHandler := handlers.NewLoginHandler(appLogger, config, createPlaylistUsecase)
	userHandler := handlers.NewUserHandler(appLogger, deleteListeningHistoryUsecase, getArtistsByUserIDUsecase)
	settingHandler := handlers.NewSettinHandler(appLogger, getSettingUsecase, editSettingUsecase)

	e.GET("/login", loginHandler.Login)
	e.GET("/callback", loginHandler.Callback)
	e.GET("/user", userHandler.GetArtists)
	e.POST("/delete", userHandler.DeleteArtists)
	e.GET("/setting", settingHandler.GetSettings)
	e.POST("/setting/save", settingHandler.EditSettings)

	return e
}

func newDB(config *config.DBConfig) (*gorp.DbMap, error) {
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

func newCron(
	logger *log.Logger,
	updateListeningHistoryUsecase *usecase.UpdateListeningHistoryUsecase,
	updatePlaylistUsecase *usecase.UpdatePlaylistUsecase,
	updateFollowingArtistsUsecase *usecase.UpdateFollowingArtistsUsecase) *cron.Cron {
	c := cron.New()
	c.AddFunc("@every 20m", func() {
		if err := updateListeningHistoryUsecase.UpdateListeningHistory(); err != nil {
			logger.Print(err)
		}
	})
	c.AddFunc("10 00 * * 5", func() {
		if err := updatePlaylistUsecase.UpdatePlaylistHistory(); err != nil {
			logger.Print(err)
		}
	})
	c.AddFunc("@monthly", func() {
		if err := updateFollowingArtistsUsecase.UpdateFollowingArtists(); err != nil {
			logger.Print(err)
		}
	})
	return c
}
