package context

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
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

func (c *AppContext) Log() *log.Logger {
	return c.Get("logger").(*zerolog.Logger)
}
