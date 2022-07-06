package utils

import (
	"strings"
)

// Error model info
// @Description Error model response information
type JError struct {
	Error string `json:"error"`
}

func NewError(err error) JError {
	jerr := JError{
		"generic error",
	}
	if err != nil {
		jerr.Error = err.Error()
	}
	return jerr
}

func NormalizeEmail(email string) string {
	normalizedEmail := strings.TrimSpace(strings.ToLower(email))
	return normalizedEmail
}

func ArrayContainsAny(arrayContainer [6]string, arrayContained []string) bool {
	for _, s := range arrayContained {
		for _, s2 := range arrayContainer {
			if s == s2 {
				return true
			}
		}
	}
	return false
}
