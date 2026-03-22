package loglint

import "regexp"

var logMethods = map[string]bool{
	"Info":    true,
	"Error":   true,
	"Warn":    true,
	"Debug":   true,
	"Trace":   true,
	"Fatal":   true,
	"Panic":   true,
	"Print":   true,
	"Println": true,
	"Printf":  true,
}

var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bpassword\b`),
	regexp.MustCompile(`(?i)\bpasswd\b`),
	regexp.MustCompile(`(?i)\btoken\b`),
	regexp.MustCompile(`(?i)api[_-]?key`),
	regexp.MustCompile(`(?i)\bsecret\b`),
	regexp.MustCompile(`(?i)\bcredential\b`),
}

var (
	nonASCIIPattern = regexp.MustCompile(`[^\x00-\x7F]`)
)
