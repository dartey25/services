package log

import "github.com/labstack/echo/v4"

type Logger interface {
	Info(msg string)
	Infof(format string, v ...interface{})
}

// LoggerMiddleware returns a middleware function that injects the logger into the echo.Context.
func LoggerMiddleware(logger Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			return next(c)
		}
	}
}
