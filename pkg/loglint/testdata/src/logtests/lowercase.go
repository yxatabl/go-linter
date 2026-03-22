package logtests

import (
	"log"
	"log/slog"
)

func GoodLowercase() {
	slog.Info("this is a lowercase message")
	log.Print("this is lowercase message")
}

func BadLowercase() {
	slog.Info("This is a bad uppercase message") // want "log message should start with a lowercase letter"
	log.Print("This is a bad uppercase message") // want "log message should start with a lowercase letter"
}
