package loglint

import (
	"strings"
	"unicode"
)

func startsWithLowercase(s string) bool {
	if len(s) == 0 {
		return true
	}

	firstRune := []rune(s)[0]
	if !unicode.IsLetter(firstRune) {
		return true
	}

	return unicode.IsLower(firstRune)
}

func isEnglishOnly(s string) bool {
	if len(s) == 0 {
		return true
	}

	return !nonASCIIPattern.MatchString(s)
}

func hasSpecialCharsOrEmojis(s string) bool {
	for _, r := range s {
		if r > 0xFFFF {
			return true
		}
	}

	if strings.Contains(s, "!!!") ||
		strings.Contains(s, "???") ||
		strings.Contains(s, "...") {
		return true
	}

	return false
}

func containsSensitiveData(s string) bool {
	sLower := strings.ToLower(s)
	for _, pattern := range defaultSensitivePatterns {
		if pattern.MatchString(sLower) {
			return true
		}
	}
	return false
}
