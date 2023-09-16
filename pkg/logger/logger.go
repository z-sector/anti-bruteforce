package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

func GetLogger(name string, jsonFormat bool) zerolog.Logger {
	var writer io.Writer = os.Stdout

	if !jsonFormat {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339Nano,
		}
	}

	return zerolog.New(writer).With().Str("logger", name).Timestamp().Caller().Logger()
}
