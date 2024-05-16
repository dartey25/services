package joker

import (
	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/pkg/eucustoms/service"
	"github.com/mdoffice/md-services/pkg/joker/handler"
)

func Register(app *echo.Echo, service *service.EuCustomService) {
	group := app.Group("joker")
	handler := handler.New(service)
	group.GET("/eori/validate", handler.HandleEoriQuery)
	group.POST("/eori/validate", handler.HandleEoriQuery)
}
