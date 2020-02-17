package zerolog

import (
	"os"

	"github.com/rs/zerolog"
)

func NewJSONLogger() *JSONLogger {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &JSONLogger{
		l: log,
	}
}

type JSONLogger struct {
	l zerolog.Logger
}

func (logger *JSONLogger) Info(message string, fields map[string]interface{}) {
	log := logger.withFields(fields)
	log.Info().Msg(message)
}

func (logger *JSONLogger) Error(err error, fields map[string]interface{}) {
	log := logger.withFields(fields)

	log.Error().Err(err).Msg(err.Error())
}

func (logger *JSONLogger) withFields(fields map[string]interface{}) zerolog.Logger {
	log := logger.l.With().Logger()
	for key, value := range fields {
		log = log.With().Interface(key, value).Logger()
	}

	return log
}
