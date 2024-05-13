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
	"github.com/mdoffice/md-services/internal/eucustoms/handler"
	"github.com/mdoffice/md-services/internal/eucustoms/service"
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
		log.Fatalf("Error connecting to Oracle: %v", err)
	}
	defer db.Close()
	// fmt.Println(1)
	// es, err := database.NewElasticClient(&cfg.Elastic)
	// if err != nil {
	// 	log.Fatalf("Error connecting to Elastic: %v", err)
	// }

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	// AllowOrigins: []string{"https://www.mdoffice.com.ua"},
	app.Use(middleware.Gzip())
	app.Pre(middleware.Rewrite(map[string]string{
		"/services/*": "/$1",
	}))
	app.Static("/static", "assets")

	euGroup := app.Group("/eucustom", middleware.Rewrite(map[string]string{
		fmt.Sprintf("/%s/*", cfg.Server.Prefix): "/$1",
	}))
	s := service.NewEuCustomService(db)
	e := handler.NewEuCustomHandler(s)
	euGroup.GET("/", e.HandleIndex)
	euGroup.GET("/aeo", e.HandleAeoTab)
	euGroup.GET("/aeo/form", e.HandleAeoForm)
	euGroup.GET("/eori", e.HandleEoriTab)
	euGroup.GET("/eori/form", e.HandleEoriForm)
	euGroup.GET("/aeo/data", e.HandleGetAeoData)
	euGroup.GET("/eori/data", e.HandleGetEoriData)
	app.GET("/joker/eori/validate", e.HandleJokerEoriData)

	// sGroup := app.Group("/sanctions")
	// ss := sanctService.NewSanctionsService(es)
	// ess := sanctHandler.NewSanctionsHandler(ss)
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
