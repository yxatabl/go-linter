package a

import (
	"log"
	"log/slog"
)

func GoodLogs() {
	slog.Info("this log starts with lowercase letter")
	slog.Error("this is error log with lowercase letter")

	log.Print("standard log check")
}

func BadLogs() {
	slog.Info("This log starts with uppercase letter")        // want "log message should start with a lowercase letter"
	slog.Error("This error log starts with uppercase letter") // want "log message should start with a lowercase letter"
}
