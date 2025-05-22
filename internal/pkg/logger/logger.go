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
	if InfoLogger == nil {
		Init()
	}
	InfoLogger.Println(message)
}

// Error logs an error message
func Error(message string, err error) {
	if ErrorLogger == nil {
		Init()
	}
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

// Fatal logs a fatal error and exits the program
func Fatal(message string, err error) {
	if ErrorLogger == nil {
		Init()
	}
	ErrorLogger.Fatalf("%s: %v\n", message, err)
}
