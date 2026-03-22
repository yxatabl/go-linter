package a

import (
	"log"
	"log/slog"
)

func AllGood() {
	// First rule tests
	slog.Info("starting server")
	slog.Error("failed to connect to database")
	slog.Warn("something happened")
	slog.Debug("debug information")

	// Second rule tests
	slog.Info("user logged in successfully")
	slog.Error("connection timeout")

	// Third rule tests
	slog.Info("server started")
	slog.Warn("warning message")

	// Fourth rule tests
	slog.Info("user authenticated")
	slog.Debug("request completed")
	slog.Error("auth failed for user")

	// Check log compability
	log.Print("standard log message")
	log.Println("log with newline")
}

// First rule tests
func BadLowercase() {
	slog.Info("Starting server")             // want "log message should start with a lowercase letter"
	slog.Error("Failed to connect")          // want "log message should start with a lowercase letter"
	slog.Warn("WARNING: something happened") // want "log message should start with a lowercase letter"
	slog.Debug("Debug message")              // want "log message should start with a lowercase letter"
}

// Second rule tests
func BadEnglish() {
	slog.Info("запуск сервера")              // want "log message must be in English only"
	slog.Error("ошибка подключения")         // want "log message must be in English only"
	slog.Info("сервер запущен successfully") // want "log message must be in English only"
	slog.Warn("предупреждение")              // want "log message must be in English only"
}

// Third rule tests
func BadSpecialChars() {
	slog.Info("server started!!!")    // want "log message should not contain special characters or emojis"
	slog.Info("connection failed...") // want "log message should not contain special characters or emojis"
	slog.Info("server started🚀")      // want "log message should not contain special characters or emojis"
	slog.Error("error!!!")            // want "log message should not contain special characters or emojis"
}

// Fourth rule tests
func BadSensitiveData() {
	slog.Info("user password: secret123")    // want "log message should not contain sensitive data"
	slog.Debug("api_key=abc123")             // want "log message should not contain sensitive data"
	slog.Info("token: eyJhbGciOiJIUzI1NiIs") // want "log message should not contain sensitive data"
	slog.Error("invalid password")           // want "log message should not contain sensitive data"
	slog.Info("user secret key")             // want "log message should not contain sensitive data"
}

// Multiple violations tests
func MultipleViolations() {
	slog.Info("Starting server!!! password=123") // want "log message should start with a lowercase letter" "log message should not contain special characters or emojis" "log message should not contain sensitive data"
}
