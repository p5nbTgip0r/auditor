package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func LoggingInitialize() {
	var writers []io.Writer

	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05 -0700"

	multi := io.MultiWriter(writers...)
	log.Logger = zerolog.New(multi).
		With().
		Timestamp().
		Logger().
		Level(zerolog.DebugLevel)
}
