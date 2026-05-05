package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(env string) zerolog.Logger {
	var output io.Writer

	if env == "development" {

		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {

		output = os.Stdout
	}

	return zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", "devix-backend").
		Logger()
}
