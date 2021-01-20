package util

import (
	"html"
	"regexp"
	"strings"
)

// ValidateName tells if a name is valid or not
func ValidateName(name string) bool {
	return len(name) > 3
}

// ValidateEmail checks if the email is valid
func ValidateEmail(email string) bool {
	reg := regexp.MustCompile(`(?:[a-z0-9!#$%&'*+/=?^_\x60{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_\x60{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`)
	return reg.Match([]byte(email))
}

// Clean removes duplicate whitespace and illegal characters
func Clean(text *string) {
	*text = strings.TrimSpace(*text)
	*text = html.EscapeString(*text)
	reg := regexp.MustCompile(`([^a-zA-Z0-9\s])|(\s{2,})`)
	*text = reg.ReplaceAllString(*text, " ")
}

// RemoveSpaces removes n duplicates and leading/trailing spaces from a string
func RemoveSpaces(password *string) {
	*password = html.EscapeString(*password)
	*password = strings.TrimSpace(*password)
	reg := regexp.MustCompile(`\s`)
	*password = reg.ReplaceAllString(*password, "")
}
