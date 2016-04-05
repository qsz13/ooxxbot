package logger

import (
	"log"
	"os"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func Info() *log.Logger {
	if infoLogger == nil {
		infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return infoLogger
}

func Warning() *log.Logger {
	if warningLogger == nil {
		warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return warningLogger
}

func Error() *log.Logger {
	if errorLogger == nil {
		errorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return errorLogger
}
