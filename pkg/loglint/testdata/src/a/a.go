package a

import (
	"log/slog"
)

func AllGood() {
	slog.Info("starting server")
	slog.Info("user logged in successfully")
	slog.Error("failed to connect to database")
}

func BadLowercase() {
	slog.Info("Starting server")             // want "log message should start with a lowercase letter"
	slog.Error("Failed to connect")          // want "log message should start with a lowercase letter"
	slog.Warn("WARNING: something happened") // want "log message should start with a lowercase letter"
}

func BadEnglish() {
	slog.Info("запуск сервера")              // want "log message must be in English only"
	slog.Error("ошибка подключения")         // want "log message must be in English only"
	slog.Info("сервер запущен successfully") // want "log message must be in English only"
}

func BadSpecialChars() {
	slog.Info("server started!!!") // want "log message should not contain special characters or emojis"
	slog.Warn("warning??")         // want "log message should not contain special characters or emojis"
}

func BadSensitiveData() {
	slog.Info("user password: secret123")    // want "log message should not contain sensitive data"
	slog.Debug("api_key=abc123")             // want "log message should not contain sensitive data"
	slog.Info("token: eyJhbGciOiJIUzI1NiIs") // want "log message should not contain sensitive data"
	slog.Error("auth failed for user")       // want "log message should not contain sensitive data"
}
