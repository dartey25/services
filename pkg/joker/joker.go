package joker

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mdoffice/md-services/pkg/eucustoms/service"
	"github.com/mdoffice/md-services/pkg/joker/handler"
)

func Register(app *echo.Echo, service *service.EuCustomService) {
	joker := app.Group("joker")
	joker.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	handler := handler.New(service)

	joker.POST("/eori/validate", handler.HandleEoriQuery)
	joker.POST("/aeo/q", handler.HandleAeoQuery)
}
