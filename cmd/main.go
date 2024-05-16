package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cfg "github.com/mdoffice/md-services/config"
	appCtx "github.com/mdoffice/md-services/internal/context"
	database "github.com/mdoffice/md-services/internal/db"
	"github.com/mdoffice/md-services/internal/eucustoms/handler"
	"github.com/mdoffice/md-services/internal/eucustoms/service"
	"github.com/mdoffice/md-services/internal/log"
)

func main() {
	logger := log.NewZeroLog()

	var cfg cfg.Config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		logger.Fatal().Err(fmt.Errorf("error reading config: %s", err.Error()))
		os.Exit(1)
	}
	db, err := database.NewOracleClient(&cfg.Database)
	if err != nil {
		logger.Fatal().Err(fmt.Errorf("error connecting to db: %s", err.Error()))
	}
	defer db.Close()

	// es, err := database.NewElasticClient(&cfg.Elastic)
	// if err != nil {
	// // 	log.Fatalf("Error connecting to Elastic: %v", err)
	// // }
	//
	app := echo.New()
	app.HideBanner = true
	app.HidePort = true

	app.Use(middleware.Recover())
	app.Use(middleware.CORS())
	app.Use(middleware.Secure())
	// AllowOrigins: []string{"https://www.mdoffice.com.ua"},
	// opts := slog.HandlerOptions{
	// 	AddSource: true,
	// 	Level:     slog.LevelInfo,
	// }

	// Use the request logger middleware with zerolog
	app.Use(middleware.RequestLoggerWithConfig(log.RequestLoggerConfig(logger)))
	app.Use(log.LoggerMiddleware(logger))

	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := appCtx.NewAppContext(c)
			return next(cc)
		}
	})

	app.Use(middleware.Gzip())
	app.Pre(middleware.Rewrite(map[string]string{
		"/services/*": "/$1",
	}))
	assetHandler := http.FileServer(rice.MustFindBox("../assets").HTTPBox())
	app.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	euGroup := app.Group("/eucustom")
	s := service.NewEuCustomService(db)
	e := handler.NewEuCustomHandler(s)
	euGroup.GET("", e.HandleIndex)
	euGroup.GET("/", e.HandleIndex)
	euGroup.GET("/aeo", e.HandleAeoTab)
	euGroup.GET("/aeo/form", e.HandleAeoForm)
	euGroup.GET("/eori", e.HandleEoriTab)
	euGroup.GET("/eori/form", e.HandleEoriForm)
	euGroup.GET("/aeo/data", e.HandleGetAeoData)
	euGroup.GET("/eori/data", e.HandleGetEoriData)
	app.GET("/joker/eori/validate", e.HandleJokerEoriData)
	app.POST("/joker/eori/validate", e.HandleJokerEoriData)

	// sGroup := app.Group("/sanctions")
	// ss := sanctionsService.NewSanctionsService(es)
	// ess := sanctionsHandler.NewSanctionsHandler(ss)
	// sGroup.GET("/parse", ess.HandleParseLegal)
	// sGroup.GET("/query", ess.HandleQueryLegal)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := app.Start(fmt.Sprintf(":%v", cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
}
