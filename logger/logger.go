package logger

import (
	"log"
	"os"
)

var (
	DebugFlag     bool
	debugLogger   *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func Debug(msg string) {
	if debugLogger == nil {
		debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	if DebugFlag {
		debugLogger.Output(2, msg)
	}
}

func Warning(msg string) {
	if warningLogger == nil {
		warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	warningLogger.Output(2, msg)
}

func Error(msg string) {
	if errorLogger == nil {
		errorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	errorLogger.Output(2, msg)
}
