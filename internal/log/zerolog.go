package log

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	zerolog.Logger
}

func NewZeroLog() *ZeroLogger {
	return &ZeroLogger{zerolog.New(os.Stdout).With().Timestamp().Logger()}
}

func (l *ZeroLogger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *ZeroLogger) Infof(format string, v ...interface{}) {
	l.Logger.Info().Msgf(format, v...)
}

// // RequestLoggerConfig creates a middleware.RequestLoggerConfig for logging requests.
func RequestLoggerConfig(l *ZeroLogger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				l.Logger.Info().
					Str("uri", v.URI).
					Int("status", v.Status).
					Msg("REQUEST")
			} else {
				l.Logger.Error().
					Err(v.Error).
					Str("uri", v.URI).
					Int("status", v.Status).
					Msg("REQUEST_ERROR")
			}
			return nil
		},
	}
}
