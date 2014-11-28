package validate

import (
	"strings"
)

func Required(field, key, label string, errors map[string]string) {
	if field == "" {
		errors[key] = label + " cannot be blank."
	}
}

func Email(field, key, label string, errors map[string]string) {
	if !strings.ContainsRune(field, '@') {
		errors[key] = label + " does not look like a valid email address."
	}
}