package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cfg "github.com/mdoffice/md-services/config"
	database "github.com/mdoffice/md-services/internal/db"
	euHandler "github.com/mdoffice/md-services/internal/eucustoms/handler"
	euService "github.com/mdoffice/md-services/internal/eucustoms/service"
	sanctHandler "github.com/mdoffice/md-services/internal/sanctions/handler"
	sanctService "github.com/mdoffice/md-services/internal/sanctions/service"
)

func main() {
	var cfg cfg.Config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config: %s", err.Error()))
		os.Exit(1)
	}
	db, err := database.NewOracleClient(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	es, err := database.NewESClient(&cfg.Elastic)
	if err != nil {
		log.Fatal(err)
	}

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	// AllowOrigins: []string{"https://www.mdoffice.com.ua"},
	app.Use(middleware.Gzip())
	app.Static("/static", "assets")

	s := euService.NewEuCustomService(db)
	e := euHandler.NewEuCustomHandler(s)
	app.GET("/", e.HandleIndex)
	app.GET("/aeo", e.HandleAeoTab)
	app.GET("/aeo/form", e.HandleAeoForm)
	app.GET("/eori", e.HandleEoriTab)
	app.GET("/eori/form", e.HandleEoriForm)
	app.GET("/aeo/data", e.HandleGetAeoData)
	app.GET("/eori/data", e.HandleGetEoriData)
	app.GET("/joker/eori/validate", e.HandleJokerEoriData)

	ss := sanctService.NewSanctionsService(es)
	ess := sanctHandler.NewSanctionsHandler(ss)
	app.GET("/sanctions/parse", ess.HandleParseLegal)
	app.GET("/sanctions/query", ess.HandleQueryLegal)

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
