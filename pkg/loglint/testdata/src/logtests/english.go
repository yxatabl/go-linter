package logtests

import (
	"log"
	"log/slog"
)

func GoodEnglish() {
	slog.Info("this is a good log message")
	log.Print("this is a good log message")
}

func BadEnglish() {
	slog.Info("лог на русском") // want "log message should not contain special characters or emojis"
	log.Print("лог на русском") // want "log message should not contain special characters or emojis"
}
