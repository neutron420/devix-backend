package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// New creates a new structured logger based on the environment.
func New(env string) zerolog.Logger {
	var output io.Writer

	if env == "development" {
		// Pretty console output for dev
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		// JSON output for production (structured, machine-parseable)
		output = os.Stdout
	}

	return zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", "devix-backend").
		Logger()
}
