package logtests

import (
	"log"
	"log/slog"
)

func GoodSensitive() {
	slog.Info("this is a good log message without sensitive information")
	log.Print("this is a good log message without sensitive information")
}

func BadSensitive() {
	slog.Info("this log message contains a password: mysecretpassword") // want "log message should not contain sensitive information"
	log.Print("this log message contains a password: mysecretpassword") // want "log message should not contain sensitive information"
}
