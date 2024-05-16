package log

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	*zerolog.Logger
}

// InitializeLogger creates and returns a new zerolog logger instance.
func NewZeroLog() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// LoggerMiddleware returns a middleware function that injects the logger into the echo.Context.
func LoggerMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", &logger)
			return next(c)
		}
	}
}

// RequestLoggerConfig creates a middleware.RequestLoggerConfig for logging requests.
func RequestLoggerConfig(logger zerolog.Logger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info().
					Str("uri", v.URI).
					Int("status", v.Status).
					Msg("REQUEST")
			} else {
				logger.Error().
					Err(v.Error).
					Str("uri", v.URI).
					Int("status", v.Status).
					Msg("REQUEST_ERROR")
			}
			return nil
		},
	}
}
