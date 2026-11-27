package sandbox

import (
	"errors"
	"regexp"
	"strings"
)

var labelKeyPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,62}$`)

// ValidateLabels accepts up to 16 pairs with 1–63 byte ASCII keys and
// values of up to 255 printable ASCII bytes. Nil and empty maps are valid.
// Credential-like keys are rejected by case-insensitive substring matching
// against token, key, secret, password, and api.
func ValidateLabels(labels map[string]string) error {
	if len(labels) > 16 {
		return errors.New("sandbox labels may contain at most 16 pairs")
	}
	for key, value := range labels {
		if !labelKeyPattern.MatchString(key) {
			return errors.New("invalid sandbox label key")
		}
		if credentialLabelKey(key) {
			return errors.New("credential-like sandbox label key")
		}
		if len(value) > 255 {
			return errors.New("sandbox label value exceeds 255 bytes")
		}
		for i := range len(value) {
			if value[i] < ' ' || value[i] > '~' {
				return errors.New("sandbox label value must be printable ASCII")
			}
		}
	}
	return nil
}

func credentialLabelKey(key string) bool {
	lower := strings.ToLower(key)
	for _, fragment := range [...]string{"token", "key", "secret", "password", "api"} {
		if strings.Contains(lower, fragment) {
			return true
		}
	}
	return false
}
