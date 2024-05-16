package context

import (
	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/log"
)

type AppContext struct {
	echo.Context
}

// NewAppContext creates a new AppContext.
func NewAppContext(c echo.Context) *AppContext {
	return &AppContext{
		Context: c,
	}
}

func (c *AppContext) Log() log.Logger {
	return c.Get("logger").(log.Logger)
}
