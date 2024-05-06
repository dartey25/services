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
	"github.com/mdoffice/md-services/db"
	"github.com/mdoffice/md-services/handlers"
	"github.com/mdoffice/md-services/model"
	"github.com/mdoffice/md-services/services"
)

func main() {
	var cfg model.Config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config: %s", err.Error()))
		os.Exit(1)
	}

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.Gzip())
	app.Static("/static", "assets")

	db, err := db.New(&cfg.Database)
	if err != nil {
		app.Logger.Fatal(err)
	}
	defer db.Close()

	s := services.NewEuCustomService(db)
	e := handlers.NewEuCustomHandler(s)
	app.GET("/", e.HandleIndex)
	app.GET("/aeo", e.HandleAeoTab)
	app.GET("/aeo/form", e.HandleAeoForm)
	app.GET("/eori", e.HandleEoriTab)
	app.GET("/eori/form", e.HandleEoriForm)
	app.GET("/aeo/data", e.HandleGetAeoData)
	app.GET("/eori/data", e.HandleGetEoriData)

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
