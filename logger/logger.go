package logger

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

// Debug is the Logger which is used for debug output
var Debug Logger = log.New(os.Stderr, "[NUKI][DEBUG]: ", log.Ldate|log.Ltime)

// Info is the Logger which is used for info output
var Info Logger = log.New(os.Stderr, "[NUKI][INFO]: ", log.Ldate|log.Ltime)
