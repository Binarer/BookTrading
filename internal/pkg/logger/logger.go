package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// Init initializes the loggers
func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs an info message
func Info(message string) {
	InfoLogger.Println(message)
}

// Error logs an error message
func Error(message string, err error) {
	ErrorLogger.Printf("%s: %v\n", message, err)
}

// Debug логирует отладочное сообщение
func Debug(msg string, fields ...interface{}) {
	log.Printf("DEBUG: %s\n", msg)
}

// Warn логирует предупреждение
func Warn(msg string, fields ...interface{}) {
	log.Printf("WARN: %s\n", msg)
}

// Fatal логирует фатальную ошибку и завершает программу
func Fatal(msg string, err error, fields ...interface{}) {
	log.Fatalf("FATAL: %s: %v\n", msg, err)
} 