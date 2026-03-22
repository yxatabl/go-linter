package logtests

import (
	"log"
	"log/slog"
)

func goodSpecial() {
	slog.Info("this is a good log message with special characters: !@#$%^&*()")
	log.Print("this is a good log message with special characters: !@#$%^&*()")
}

func badSpecial() {
	slog.Info("this log message contains an emoji: 😊") // want "log message should not contain special characters or emojis"
	log.Print("this log message contains an emoji: 😊") // want "log message should not contain special characters or emojis"
}
