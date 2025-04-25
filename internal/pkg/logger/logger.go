package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger инициализирует логгер с заданными настройками
func InitLogger(level zerolog.Level, format string) {
	// Настройка временной зоны для логов
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(level)

	// Настройка формата вывода
	if format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		log.Logger = log.Output(os.Stdout)
	}
}

// Info логирует информационное сообщение
func Info(msg string, fields ...interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

// Error логирует сообщение об ошибке
func Error(msg string, err error, fields ...interface{}) {
	log.Error().Err(err).Fields(fields).Msg(msg)
}

// Debug логирует отладочное сообщение
func Debug(msg string, fields ...interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

// Warn логирует предупреждение
func Warn(msg string, fields ...interface{}) {
	log.Warn().Fields(fields).Msg(msg)
}

// Fatal логирует фатальную ошибку и завершает программу
func Fatal(msg string, err error, fields ...interface{}) {
	log.Fatal().Err(err).Fields(fields).Msg(msg)
} 