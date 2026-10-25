package main

import (
	"errors"
	"strings"
)

func redactedError(message string, secrets ...string) error {
	return errors.New(redactSecrets(message, secrets...))
}

func redactSecrets(message string, secrets ...string) string {
	for _, secret := range secrets {
		if secret != "" {
			message = strings.ReplaceAll(message, secret, "[REDACTED]")
		}
	}
	return message
}
