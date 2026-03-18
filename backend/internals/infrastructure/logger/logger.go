package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func GetLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger
}
