package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(env string, level string) zerolog.Logger {
	var output io.Writer

	if env == "development" {

		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {

		output = os.Stdout
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	return zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Str("service", "devix-backend").
		Logger()
}
